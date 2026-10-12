package kv

import (
	"path/filepath"
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

func TestPersistenceSaveAndRecover(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test_data.json")

	// Phase 1: Initialize store and insert data
	store1, err := NewPersistentStore(tempFile)
	if err != nil {
		t.Fatalf("failed to create persistent store: %v", err)
	}

	store1.Set("k1", "v1")
	store1.Set("k2", "v2")

	// Phase 2: Simulate restart by creating a new store instance with same file
	store2, err := NewPersistentStore(tempFile)
	if err != nil {
		t.Fatalf("failed to recover store after restart: %v", err)
	}

	val1, ok1 := store2.Get("k1")
	if !ok1 || val1 != "v1" {
		t.Fatalf("expected k1=v1 after restart, got val=%s, ok=%v", val1, ok1)
	}

	val2, ok2 := store2.Get("k2")
	if !ok2 || val2 != "v2" {
		t.Fatalf("expected k2=v2 after restart, got val=%s, ok=%v", val2, ok2)
	}

	// Phase 3: Delete a key and restart again
	store2.Delete("k1")

	store3, err := NewPersistentStore(tempFile)
	if err != nil {
		t.Fatalf("failed to recover store after 2nd restart: %v", err)
	}

	_, exists := store3.Get("k1")
	if exists {
		t.Fatal("expected k1 to be deleted across restart")
	}

	val2Again, ok2Again := store3.Get("k2")
	if !ok2Again || val2Again != "v2" {
		t.Fatalf("expected k2=v2 to persist, got val=%s, ok=%v", val2Again, ok2Again)
	}
}
