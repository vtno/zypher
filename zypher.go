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

type CipherFactory struct{}

func NewCipherFactory() *CipherFactory {
	return &CipherFactory{}
}

func (cf *CipherFactory) NewCipher(key string) crypto.Cipher {
	return NewCipher(key)
}

type Cipher struct {
	key []byte
}

func NewCipher(key string) *Cipher {
	return &Cipher{
		key: []byte(key),
	}
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
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

func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
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
