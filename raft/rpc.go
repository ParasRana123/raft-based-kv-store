package raft

// RequestVoteArgs is the RPC payload for requesting a vote during election.
type RequestVoteArgs struct {
	Term         int    `json:"term"`
	CandidateID  string `json:"candidate_id"`
	LastLogIndex int    `json:"last_log_index"`
	LastLogTerm  int    `json:"last_log_term"`
}

// RequestVoteReply is the RPC response to a vote request.
type RequestVoteReply struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"vote_granted"`
}

// AppendEntriesArgs is the RPC payload for log replication and heartbeats.
type AppendEntriesArgs struct {
	Term         int        `json:"term"`
	LeaderID     string     `json:"leader_id"`
	PrevLogIndex int        `json:"prev_log_index"`
	PrevLogTerm  int        `json:"prev_log_term"`
	Entries      []LogEntry `json:"entries"`
	LeaderCommit int        `json:"leader_commit"`
}

// AppendEntriesReply is the RPC response to an AppendEntries call.
type AppendEntriesReply struct {
	Term    int  `json:"term"`
	Success bool `json:"success"`
}
