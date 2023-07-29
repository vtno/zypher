package store_test

import (
	"log"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/vtno/zypher/internal/server/store"
)

func MustSet(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error setting value: %v", err)
	}
}

type GetResult struct {
	value string
	err   error
}

func MustGetExpected(t *testing.T, expected string, gr *GetResult) {
	if gr.err != nil {
		t.Errorf("error getting value: %v", gr.err)
	}
	if gr.value != expected {
		t.Errorf("expected value to be %s, got %s", expected, gr.value)
	}
}

func TestBBoltStore_GetSet(t *testing.T) {
	t.Run("successfully sets a value with the provided key", func(t *testing.T) {
		store, err := store.NewBBoltStore("test.db")
		if err != nil {
			t.Errorf("error creating store: %v", err)
		}
		defer func() {
			store.Close()
			err := os.Remove("test.db")
			if err != nil {
				log.Fatalf("error removing test.db: %v", err)
			}
		}()

		MustSet(t, store.Set("prod#somekey", "somevalue"))
		value, err := store.Get("prod#somekey")
		MustGetExpected(t, "somevalue", &GetResult{
			value: value,
			err:   err,
		})
		MustSet(t, store.Set("prod#somekey", "newvalue"))
		value, err = store.Get("prod#somekey")
		MustGetExpected(t, "newvalue", &GetResult{
			value: value,
			err:   err,
		})
	})
}

func TestBBoltStore_List(t *testing.T) {
	type test struct {
		name     string
		expected []string
		prefix   *string
	}
	prodPrefix := "prod"
	tests := []test{
		{
			name:     "successfully lists all the keys in the store when prefix not provided",
			expected: []string{"prod#key1", "prod#key2", "dev#key1"},
			prefix:   nil,
		},
		{
			name:     "successfully lists all the keys in the store when prefix provided",
			expected: []string{"prod#key1", "prod#key2"},
			prefix:   &prodPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := store.NewBBoltStore("test.db")
			if err != nil {
				t.Errorf("error creating store: %v", err)
			}
			defer func() {
				store.Close()
				err := os.Remove("test.db")
				if err != nil {
					log.Fatalf("error removing test.db: %v", err)
				}
			}()

			MustSet(t, store.Set("prod#key1", "somevalue1"))
			MustSet(t, store.Set("prod#key2", "somevalue2"))
			MustSet(t, store.Set("dev#key1", "somevalue3"))

			values, err := store.List(tt.prefix)
			if err != nil {
				t.Errorf("error listing keys: %v", err)
			}
			if len(values) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(values))
			}
			sort.Strings(values)
			sort.Strings(tt.expected)
			if !reflect.DeepEqual(values, tt.expected) {
				t.Errorf("expected %v keys, got %v", tt.expected, values)
			}
		})
	}
}
