package kv

// KVStore defines the operations supported by the key-value storage layer.
type KVStore interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// Store is an in-memory implementation of the KVStore interface.
type Store struct {
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

func (s *Store) Set(key string, value string) {
	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	value, exists := s.data[key]
	return value, exists
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}
