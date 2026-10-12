package kv

import (
	"sync"
	"testing"
)

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

	store.Delete("non_existent")

	_, exists := store.Get("non_existent")
	if exists {
		t.Fatal("expected key not to exist")
	}
}

func TestKVStoreInterface(t *testing.T) {
	var store KVStore = NewStore()

	store.Set("interface_key", "value")
	val, ok := store.Get("interface_key")
	if !ok || val != "value" {
		t.Fatalf("expected 'value', got '%s'", val)
	}

	store.Delete("interface_key")
	_, ok = store.Get("interface_key")
	if ok {
		t.Fatal("expected key to be deleted via KVStore interface")
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewStore()
	var wg sync.WaitGroup

	numGoroutines := 50

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Set(string(rune('a'+id%26)), "val")
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Get(string(rune('a' + id%26)))
		}(i)
	}

	// Concurrent deletes
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Delete(string(rune('a' + id%26)))
		}(i)
	}

	wg.Wait()
}


