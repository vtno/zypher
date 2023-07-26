package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/vtno/zypher/internal/server/store"
	"golang.org/x/crypto/ssh"
)

const defaultAuthBucket = "auth"

type Auth struct {
	store      store.Store
	pkProvider PubKeyProvider
	rootPubKey *rsa.PublicKey
}

type PubKeyProvider interface {
	Get() (*rsa.PublicKey, error)
}

type AuthOption func(*Auth)

func WithPubKeyProvider(p PubKeyProvider) AuthOption {
	return func(a *Auth) {
		a.pkProvider = p
	}
}

func parseTokenAndSigFromAuthHeader(authHeader string) (string, []byte) {
	s := strings.Split(authHeader, " ")
	if len(s) == 2 {
		t := strings.Split(s[1], ":")
		if len(t) == 2 {
			sig, err := base64.StdEncoding.DecodeString(t[1])
			if err != nil {
				fmt.Printf("%e\n", fmt.Errorf("error decoding signature: %w", err))
				return "", nil
			}
			return t[0], sig
		}
	}
	return "", nil
}

func NewAuth(store store.Store, opts ...AuthOption) (*Auth, error) {
	a := &Auth{store: store}
	for _, opt := range opts {
		opt(a)
	}
	if a.pkProvider == nil {
		return nil, fmt.Errorf("no pubkeyprovider is provided")
	}
	pk, err := a.pkProvider.Get()
	if err != nil {
		return nil, fmt.Errorf("error getting public key: %w", err)
	}
	a.rootPubKey = pk
	return a, nil
}

// Authenticate extracts the signed payload from the Authorization header
// then attempts to verify the signature using the public key stored in the store
func (a *Auth) Authenticate(authHeader string) bool {
	if authHeader == "" {
		return false
	}
	token, sig := parseTokenAndSigFromAuthHeader(authHeader)
	hashed := sha256.Sum256([]byte(token))
	err := rsa.VerifyPKCS1v15(a.rootPubKey, crypto.SHA256, hashed[:], sig)

	if err != nil {
		fmt.Printf("%e\n", fmt.Errorf("error verifying token: %w", err))
		return false
	}
	return true
}

func validate(base64PubKey string) error {
	pub, err := base64.StdEncoding.DecodeString(base64PubKey)
	if err != nil {
		return fmt.Errorf("error decoding public key: %w", err)
	}
	_, err = ssh.ParsePublicKey(pub)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}
	return nil
}
