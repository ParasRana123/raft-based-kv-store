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
