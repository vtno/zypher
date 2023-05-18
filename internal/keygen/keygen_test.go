package keygen_test

import (
	"testing"

	"github.com/vtno/zypher/internal/keygen"
)

func Test_GenerateKey(t *testing.T) {
	t.Parallel()

	key, err := keygen.GenerateKey()
	if err != nil {
		t.Fatalf("error generating key: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("expected key length of 32, got %d", len(key))
	}
}
