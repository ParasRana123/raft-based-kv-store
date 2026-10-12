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
