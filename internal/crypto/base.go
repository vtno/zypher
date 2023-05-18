package crypto

import (
	"flag"
	"fmt"
	"os"

	"github.com/vtno/zypher/internal/config"
)

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type CipherFactory interface {
	NewCipher(string) Cipher
}

type FileReaderWriter interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, os.FileMode) error
}

type BaseCmd struct {
	frw FileReaderWriter
	fs  *flag.FlagSet
	cfg *config.Config
	cf  CipherFactory
	ci  Cipher
}

func WithFileReaderWriter(frw FileReaderWriter) func(*BaseCmd) {
	return func(c *BaseCmd) {
		c.frw = frw
	}
}

// init parse the flags and set the config struct
// it also read the key from the zypher.key file or env variable
// key file location is overridden if the config is parsed from flags
func (b *BaseCmd) init(args []string) error {
	err := b.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing flag from args: %w", err)
	}

	if len(b.fs.Args()) > 0 {
		b.cfg.Input = b.fs.Args()[0]
	}

	if b.cfg.Key == "" {
		key, _ := b.frw.ReadFile(b.cfg.KeyFile)
		fmt.Println("key: ", key)
		if key != nil {
			b.cfg.Key = string(key)
		}
	}

	if b.cfg.Key == "" {
		key, found := os.LookupEnv("ZYPHER_KEY")
		if !found {
			return fmt.Errorf("no key provided: ")
		}
		b.cfg.Key = key
	}

	b.ci = b.cf.NewCipher(b.cfg.Key)
	return nil
}
