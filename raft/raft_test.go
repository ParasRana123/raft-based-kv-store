package raft

import (
	"testing"
)

func TestNewRaftNode(t *testing.T) {
	applyCh := make(chan ApplyMsg, 16)
	node := NewRaftNode("node-1", []string{"node-2", "node-3"}, applyCh)

	term, isLeader, role := node.GetState()
	if term != 0 {
		t.Fatalf("expected initial term 0, got %d", term)
	}
	if isLeader {
		t.Fatal("expected new node not to be leader")
	}
	if role != Follower {
		t.Fatalf("expected role Follower, got %s", role)
	}
	if node.ID() != "node-1" {
		t.Fatalf("expected ID node-1, got %s", node.ID())
	}
}

func TestRequestVote(t *testing.T) {
	applyCh := make(chan ApplyMsg, 16)
	node := NewRaftNode("node-1", []string{"node-2", "node-3"}, applyCh)

	// 1. Candidate with lower term should be rejected
	var reply RequestVoteReply
	node.RequestVote(&RequestVoteArgs{
		Term:         -1,
		CandidateID:  "node-2",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}, &reply)

	if reply.VoteGranted {
		t.Fatal("expected vote to be rejected for stale term")
	}

	// 2. Candidate with higher term should be granted vote
	node.RequestVote(&RequestVoteArgs{
		Term:         1,
		CandidateID:  "node-2",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}, &reply)

	if !reply.VoteGranted {
		t.Fatal("expected vote to be granted for higher term")
	}

	term, _, _ := node.GetState()
	if term != 1 {
		t.Fatalf("expected node term to be updated to 1, got %d", term)
	}

	// 3. Second candidate in same term should be rejected
	var reply2 RequestVoteReply
	node.RequestVote(&RequestVoteArgs{
		Term:         1,
		CandidateID:  "node-3",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}, &reply2)

	if reply2.VoteGranted {
		t.Fatal("expected second candidate in same term to be rejected")
	}
}
