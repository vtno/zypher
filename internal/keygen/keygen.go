package keygen

import (
	"crypto/rand"
	"fmt"

	"github.com/vtno/zypher/internal/file"
)

type KeyGenCmd struct {
	frw *file.FileReaderWriter
}

const (
	HelpMsg = `Usage: zypher keygen
	generates a AES-256 and save it to zypher.key file
	`
	SynopsisMsg = "generates a new key"
)

func WithFileReaderWriter(frw *file.FileReaderWriter) func(*KeyGenCmd) {
	return func(k *KeyGenCmd) {
		k.frw = frw
	}
}

func NewKeyGenCmd(opts ...func(*KeyGenCmd)) *KeyGenCmd {
	keyGenCmd := &KeyGenCmd{
		frw: file.NewFileReaderWriter(),
	}

	for _, opt := range opts {
		opt(keyGenCmd)
	}

	return keyGenCmd
}

func (k *KeyGenCmd) Help() string {
	return HelpMsg
}

func (k *KeyGenCmd) Synopsis() string {
	return SynopsisMsg
}

func (k *KeyGenCmd) Run(args []string) int {
	key, err := GenerateKey()
	if err != nil {
		fmt.Printf("error generating key: %v\n", err)
		return 1
	}

	if err := k.saveKey(key); err != nil {
		fmt.Printf("error saving key: %v\n", err)
		return 1
	}
	return 0
}

func (k *KeyGenCmd) saveKey(key string) error {
	return k.frw.WriteFile("zypher.key", []byte(key), 0600)
}

// GenerateKey generates a random string key of len 32 for AES-256 encryption
func GenerateKey() (string, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", key), nil
}
