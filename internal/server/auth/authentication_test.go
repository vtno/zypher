package auth_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/vtno/zypher/internal/server/auth"
	"github.com/vtno/zypher/internal/server/store"
	"go.uber.org/mock/gomock"
)

func createKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("error generating RSA key pair: %v", err)
	}
	return privKey, &privKey.PublicKey
}

func TestAuth_AuthenticateRoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProvider := auth.NewMockPubKeyProvider(ctrl)
	mStore := store.NewMockStore(ctrl)
	privKey, pub := createKeys(t)
	mProvider.EXPECT().Get().Return(pub, nil).AnyTimes()
	
	hashToken := sha256.Sum256([]byte("token"))
	sig, err := privKey.Sign(rand.Reader, hashToken[:], crypto.SHA256)
	if err != nil {
		t.Errorf("error signing token: %v", err)
	}

	type test struct {
		name           string
		args           []string
		expectedResult bool
	}

	tests := []test{
		{
			name: "should return false when authHeader is empty",
			args:           []string{""},
			expectedResult: false,
		},
		{
			name: "should return false when signature is invalid",
			args:           []string{"Bearer invalid:signature"},
			expectedResult: false,
		},
		{
			name: "should return true when signature is valid",
			args:           []string{"Bearer token:" + base64.StdEncoding.EncodeToString(sig)},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := auth.NewAuth(mStore, auth.WithPubKeyProvider(mProvider))
			if err != nil {
				t.Errorf("error initializing auth: %v", err)
			}
			result := auth.AuthenticateRoot(tt.args[0])
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
