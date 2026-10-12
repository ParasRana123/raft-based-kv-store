package kv

import (
	"encoding/json"
	"os"
	"sync"
)

// KVStore defines the operations supported by the key-value storage layer.
type KVStore interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// Store is a thread-safe in-memory and persistent implementation of the KVStore interface.
type Store struct {
	mu       sync.RWMutex
	data     map[string]string
	filePath string
}

// Compile-time check to ensure Store implements KVStore.
var _ KVStore = (*Store)(nil)

// NewStore initializes and returns a new Store instance.
// If an optional filePath is provided, persistence is enabled and existing data is loaded.
func NewStore(filePath ...string) *Store {
	s := &Store{
		data: make(map[string]string),
	}
	if len(filePath) > 0 && filePath[0] != "" {
		s.filePath = filePath[0]
		_ = s.loadFromFile()
	}
	return s
}

// NewPersistentStore creates a new Store backed by the specified file, loading any existing state.
func NewPersistentStore(filePath string) (*Store, error) {
	s := &Store{
		data:     make(map[string]string),
		filePath: filePath,
	}
	if err := s.loadFromFile(); err != nil {
		return nil, err
	}
	return s, nil
}

// Set inserts or updates a key-value pair in a thread-safe manner and flushes to disk.
func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	_ = s.saveToFile()
}

// Get retrieves the value for a given key with shared read access.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Delete removes a key-value pair in a thread-safe manner and updates disk.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	_ = s.saveToFile()
}

// loadFromFile reads and unmarshals the JSON data from disk if the file exists.
func (s *Store) loadFromFile() error {
	if s.filePath == "" {
		return nil
	}

	bytes, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &s.data)
}

// saveToFile persists the current in-memory map to disk as JSON.
// Must be called while holding s.mu (write lock).
func (s *Store) saveToFile() error {
	if s.filePath == "" {
		return nil
	}

	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, bytes, 0644)
}
