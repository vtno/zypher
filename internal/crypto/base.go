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

type FileReader func(string) ([]byte, error)

type BaseCmd struct {
	fr  FileReader
	fs  *flag.FlagSet
	cfg *config.Config
	cf  CipherFactory
	ci  Cipher
}

func WithFileReader(fr FileReader) func(*BaseCmd) {
	return func(c *BaseCmd) {
		c.fr = fr
	}
}

func (b *BaseCmd) init(args []string) error {
	err := b.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing flag from args: %w", err)
	}
	if len(b.fs.Args()) > 0 {
		b.cfg.Input = b.fs.Args()[0]
	}
	if b.cfg.Key == "" {
		key, found := os.LookupEnv("ZYPHER_KEY")
		if !found {
			return fmt.Errorf("no key provided")
		}
		b.cfg.Key = key
	}

	b.ci = b.cf.NewCipher(b.cfg.Key)
	return nil
}
