package store

// Store is an interface for storing and retrieving key-value pairs from different store implementations.
type Store interface {
	// Get retrieves the value associated with the given key.
	Get(key string) (string, error)
	// Set stores the given value and associates it with the given key.
	Set(key, value string) error
	// Delete removes the value associated with the given key.
	Delete(key string) error
	// List returns all the keys in the store.
	// It optionally takes a prefix to filter the keys by.
	List(prefix *string) ([]string, error)
}
