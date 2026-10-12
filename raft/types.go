package raft

// NodeRole represents the current consensus role of a Raft node.
type NodeRole string

const (
	Follower  NodeRole = "Follower"
	Candidate NodeRole = "Candidate"
	Leader    NodeRole = "Leader"
)

// LogEntry represents an individual log entry replicated across the cluster.
type LogEntry struct {
	Index int    `json:"index"`
	Term  int    `json:"term"`
	Data  string `json:"data"`
}

// ApplyMsg represents a committed command sent to the state machine.
type ApplyMsg struct {
	CommandValid bool
	Command      string
	CommandIndex int
	CommandTerm  int
}
