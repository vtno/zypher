// zypher is a package that provides Advanced Encryption Standard (AES) encryption and decryption.
// It is a thin wrapper around the standard library crypto/aes package which make it easier to use.
// It uses GCM mode by default.
package zypher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/vtno/zypher/internal/crypto"
)

// CipherFactory is a factory for creating a Cipher struct.
// It is needed in because the Cipher struct cannot be initialized on
// package initialization because it needs a key to be initialized.
// The key is provided by the user via the command line and the factory is used on each command.
// For usage as a library, you can use Cipher directly.
type CipherFactory struct{}

// NewCipherFactory returns a new CipherFactory.
func NewCipherFactory() *CipherFactory {
	return &CipherFactory{}
}

// NewCipher returns a new Cipher struct initialized with provided key.
func (cf *CipherFactory) NewCipher(key string) crypto.Cipher {
	return NewCipher(key)
}

// Cipher is a struct that holds the key used for encryption and decryption.
type Cipher struct {
	key []byte
}

// NewCipher returns a new Cipher struct initialized with provided key.
func NewCipher(key string) *Cipher {
	return &Cipher{
		key: []byte(key),
	}
}

// Encrypt encrypts the provided plaintext and returns the ciphertext or err.
func (c *Cipher) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	ci, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("error creating aes.Cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(ci)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher.GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("error randomizing nounce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts the provided ciphertext and returns the plaintext or err.
func (c *Cipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	ci, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("error creating aes.Cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(ci)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher.GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
