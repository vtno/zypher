package encrypt

import (
	"flag"
	"fmt"

	"github.com/vtno/zypher"
	"github.com/vtno/zypher/internal/config"
)

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type CipherFactory interface {
	NewCipher(string) Cipher
}

const (
	HelpMsg = `Usage: zypher encrypt [options] <input-value>
available options:
	-k, --key=<key>			key to encrypt/decrypt
	-f, --file=<path-to-file>	input file to be encrypted
	-o, --out=<path-to-file>	output file to be created
`
	SynopsisMsg = "encrypts input value or file with the provided key and prints the encrypted value to stdout or create a file"
)

type EncryptCmd struct {
	fs  *flag.FlagSet
	cfg *config.Config
	cf  CipherFactory
	ci  Cipher
}

func NewEncryptCmd(cf CipherFactory) *EncryptCmd {
	fs := flag.NewFlagSet("encrypt", flag.ContinueOnError)
	cfg := &config.Config{}
	fs.StringVar(&cfg.Key, "key", "", "key to encrypt/decrypt")
	fs.StringVar(&cfg.Key, "k", "", "key to encrypt/decrypt (shorthand)")
	fs.StringVar(&cfg.OutFile, "out", "", "output file to be created")
	fs.StringVar(&cfg.OutFile, "o", "", "output file to be created (shorthand)")
	fs.StringVar(&cfg.InputFile, "file", "", "input file to be encrypted")
	fs.StringVar(&cfg.InputFile, "f", "", "input file to be encrypted (shorthand)")

	return &EncryptCmd{
		cfg: cfg,
		fs:  fs,
		cf:  cf,
	}
}

func (e *EncryptCmd) Help() string {
	return HelpMsg
}

func (e *EncryptCmd) Synopsis() string {
	return SynopsisMsg
}

func (e *EncryptCmd) Init(args []string) error {
	err := e.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing flag from args: %w", err)
	}
	e.ci = zypher.NewCipher(e.cfg.Key)
	return nil
}

func (e *EncryptCmd) Run(args []string) int {
	e.Init(args)
	if len(e.fs.Args()) == 1 {
		encrypted, err := e.ci.Encrypt([]byte(e.fs.Args()[0]))
		if err != nil {
			fmt.Printf("error encrypting: %v\n", err)
			return 1
		}
		fmt.Printf(string(encrypted))
	}
	return 0
}
