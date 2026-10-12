package kv

import "sync"

// KVStore defines the operations supported by the key-value storage layer.
type KVStore interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// Store is a thread-safe in-memory implementation of the KVStore interface.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// Compile-time check to ensure Store implements KVStore.
var _ KVStore = (*Store)(nil)

// NewStore initializes and returns a new Store instance.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Set inserts or updates a key-value pair in a thread-safe manner.
func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

// Get retrieves the value for a given key with shared read access.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Delete removes a key-value pair in a thread-safe manner.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
