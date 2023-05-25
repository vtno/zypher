package store

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	defaultBucket = []byte("zypher")
)

type BBoltStore struct {
	DB *bolt.DB
}

// NewBBoltStore creates a new BBoltStore.
// it opens the db file and creates a bucket if it doesn't exist.
func NewBBoltStore(dbFilePath string) (*BBoltStore, error) {
	db, err := bolt.Open(dbFilePath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening db file for bolt: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error creating bbolt bucket: %w", err)
	}

	return &BBoltStore{
		DB: db,
	}, nil
}

// Close is a delegate function that closes the underlying db file.
func (b *BBoltStore) Close() error {
	return b.DB.Close()
}

// Get retrieves the value associated with the given key.
func (b *BBoltStore) Get(key string) (string, error) {
	var value []byte
	err := b.DB.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(defaultBucket).Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error getting value from %s bucket: %w", defaultBucket, err)
	}
	return string(value), nil
}

// Set stores the given value and associates it with the given key.
func (b *BBoltStore) Set(key, value string) error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put([]byte(key), []byte(value))
	})
	if err != nil {
		return fmt.Errorf("error setting value in %s bucket: %w", defaultBucket, err)
	}
	return nil
}

// Delete removes the value associated with the given key.
func (b *BBoltStore) Delete(key string) error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Delete([]byte(key))
	})
	if err != nil {
		return fmt.Errorf("error deleting value from %s bucket: %w", defaultBucket, err)
	}
	return nil
}

// List returns all the keys in the store.
// It optionally takes a prefix to filter the keys by.
func (b *BBoltStore) List(prefix *string) ([]string, error) {
	var keys []string
	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		c := bucket.Cursor()

		if prefix != nil {
			prefixBytes := []byte(*prefix)
			for k, _ := c.Seek(prefixBytes); k != nil && len(k) > 0 && string(k[:len(prefixBytes)]) == *prefix; k, _ = c.Next() {
				keys = append(keys, string(k))
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				keys = append(keys, string(k))
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error listing keys from %s bucket: %w", defaultBucket, err)
	}
	return keys, nil
}
