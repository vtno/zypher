package encrypt

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
	fr  FileReader
}

func WithFileReader(fr FileReader) func(*EncryptCmd) {
	return func(e *EncryptCmd) {
		e.fr = fr
	}
}

func NewEncryptCmd(cf CipherFactory, opts ...func(*EncryptCmd)) *EncryptCmd {
	fs := flag.NewFlagSet("encrypt", flag.ContinueOnError)
	cfg := &config.Config{}
	fs.StringVar(&cfg.Key, "key", "", "key to encrypt/decrypt")
	fs.StringVar(&cfg.Key, "k", "", "key to encrypt/decrypt (shorthand)")
	fs.StringVar(&cfg.OutFile, "out", "", "output file to be created")
	fs.StringVar(&cfg.OutFile, "o", "", "output file to be created (shorthand)")
	fs.StringVar(&cfg.InputFile, "file", "", "input file to be encrypted")
	fs.StringVar(&cfg.InputFile, "f", "", "input file to be encrypted (shorthand)")

	e := &EncryptCmd{
		cfg: cfg,
		fs:  fs,
		cf:  cf,
		fr:  os.ReadFile,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *EncryptCmd) init(args []string) error {
	err := e.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing flag from args: %w", err)
	}
	if len(e.fs.Args()) > 0 {
		e.cfg.Input = e.fs.Args()[0]
	}
	e.ci = e.cf.NewCipher(e.cfg.Key)
	return nil
}

func (e *EncryptCmd) Help() string {
	return HelpMsg
}

func (e *EncryptCmd) Synopsis() string {
	return SynopsisMsg
}

func (e *EncryptCmd) Run(args []string) int {
	if err := e.init(args); err != nil {
		fmt.Printf("error initializing encrypt cmd: %v", err)
	}

	var (
		input []byte
		err   error
	)

	if e.cfg.Input != "" {
		input = []byte(e.cfg.Input)
	}

	if e.cfg.InputFile != "" {
		input, err = e.fr(e.cfg.InputFile)
		if err != nil {
			fmt.Printf("error reading input file: %v\n", err)
			return 1
		}
	}

	encrypted, err := e.ci.Encrypt(input)
	if err != nil {
		fmt.Printf("error encrypting: %v\n", err)
		return 1
	}
	fmt.Println(string(encrypted))
	return 0
}
