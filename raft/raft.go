package raft

import (
	"sync"
)

// RaftNode represents a single consensus node in the Raft cluster.
type RaftNode struct {
	mu sync.RWMutex

	id    string
	peers []string

	// Persistent state on all servers
	currentTerm int
	votedFor    string
	log         []LogEntry

	// Volatile state on all servers
	commitIndex int
	lastApplied int
	role        NodeRole

	// Volatile state on leaders
	nextIndex  map[string]int
	matchIndex map[string]int

	applyCh chan ApplyMsg
}

// NewRaftNode creates and initializes a new Raft consensus node.
func NewRaftNode(id string, peers []string, applyCh chan ApplyMsg) *RaftNode {
	rn := &RaftNode{
		id:          id,
		peers:       peers,
		currentTerm: 0,
		votedFor:    "",
		log:         make([]LogEntry, 0),
		commitIndex: 0,
		lastApplied: 0,
		role:        Follower,
		nextIndex:   make(map[string]int),
		matchIndex:  make(map[string]int),
		applyCh:     applyCh,
	}

	return rn
}

// ID returns the identifier of this Raft node.
func (rn *RaftNode) ID() string {
	return rn.id
}

// GetState returns currentTerm, whether this node believes it is the leader, and current role.
func (rn *RaftNode) GetState() (int, bool, NodeRole) {
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	return rn.currentTerm, rn.role == Leader, rn.role
}

// RequestVote handles incoming vote requests from candidate nodes.
func (rn *RaftNode) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// 1. Reply false if term < currentTerm
	if args.Term < rn.currentTerm {
		reply.Term = rn.currentTerm
		reply.VoteGranted = false
		return
	}

	// If term is higher, step down to Follower and update term
	if args.Term > rn.currentTerm {
		rn.currentTerm = args.Term
		rn.role = Follower
		rn.votedFor = ""
	}

	reply.Term = rn.currentTerm

	// 2. Check if already voted in this term
	canVote := rn.votedFor == "" || rn.votedFor == args.CandidateID

	// Check candidate's log is up-to-date
	lastLogIndex := len(rn.log)
	lastLogTerm := 0
	if lastLogIndex > 0 {
		lastLogTerm = rn.log[lastLogIndex-1].Term
	}

	logOk := args.LastLogTerm > lastLogTerm ||
		(args.LastLogTerm == lastLogTerm && args.LastLogIndex >= lastLogIndex)

	if canVote && logOk {
		rn.votedFor = args.CandidateID
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
}

// AppendEntries handles log replication and heartbeat RPCs from the leader.
func (rn *RaftNode) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// 1. Reply false if term < currentTerm
	if args.Term < rn.currentTerm {
		reply.Term = rn.currentTerm
		reply.Success = false
		return
	}

	// If term is greater or equal, acknowledge leader and step down to Follower
	if args.Term > rn.currentTerm || rn.role != Follower {
		rn.currentTerm = args.Term
		rn.role = Follower
		rn.votedFor = ""
	}

	reply.Term = rn.currentTerm

	// 2. Check log consistency at prevLogIndex
	if args.PrevLogIndex > 0 {
		if args.PrevLogIndex > len(rn.log) {
			reply.Success = false
			return
		}
		if rn.log[args.PrevLogIndex-1].Term != args.PrevLogTerm {
			reply.Success = false
			return
		}
	}

	// 3. Append any new entries not already in the log
	for i, entry := range args.Entries {
		index := args.PrevLogIndex + i + 1
		if index <= len(rn.log) {
			if rn.log[index-1].Term != entry.Term {
				rn.log = rn.log[:index-1]
				rn.log = append(rn.log, entry)
			}
		} else {
			rn.log = append(rn.log, entry)
		}
	}

	// 4. Update commit index
	if args.LeaderCommit > rn.commitIndex {
		lastNewIndex := args.PrevLogIndex + len(args.Entries)
		if args.LeaderCommit < lastNewIndex {
			rn.commitIndex = args.LeaderCommit
		} else {
			rn.commitIndex = lastNewIndex
		}
	}

	reply.Success = true
}
