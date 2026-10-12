package kv

import "testing"

func TestSetAndGet(t *testing.T) {
	store := NewStore()

	store.Set("name", "Paras")

	value, exists := store.Get("name")
	if !exists {
		t.Fatal("expected key to exist")
	}

	if value != "Paras" {
		t.Fatalf("expected Paras, got %s", value)
	}
}

func TestUpdate(t *testing.T) {
	store := NewStore()

	store.Set("key", "val1")
	store.Set("key", "val2")

	value, exists := store.Get("key")
	if !exists {
		t.Fatal("expected key to exist")
	}

	if value != "val2" {
		t.Fatalf("expected val2, got %s", value)
	}
}

func TestGetMissingKey(t *testing.T) {
	store := NewStore()

	_, exists := store.Get("missing")
	if exists {
		t.Fatal("expected key not to exist")
	}
}

func TestDelete(t *testing.T) {
	store := NewStore()

	store.Set("name", "Paras")
	store.Delete("name")

	_, exists := store.Get("name")
	if exists {
		t.Fatal("expected key to be deleted")
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	store := NewStore()

	// Deleting a non-existent key should not panic
	store.Delete("non_existent")

	_, exists := store.Get("non_existent")
	if exists {
		t.Fatal("expected key not to exist")
	}
}
