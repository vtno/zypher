package provider

import (
	"crypto/rsa"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type PubKeyProvider struct {
	path string
}

func NewPubKeyProvider(path string) *PubKeyProvider {
	return &PubKeyProvider{
		path: path,
	}
}

func (p *PubKeyProvider) Get() (*rsa.PublicKey, error) {
	fmt.Printf("loading root public key from %s\n", p.path)
	pub, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("error reading public key at %s: %w", p.path, err)

	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pub)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}
	return pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey), nil
}
