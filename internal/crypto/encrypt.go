package crypto

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/vtno/zypher/internal/config"
	"github.com/vtno/zypher/internal/file"
)

const (
	HelpMsg = `Usage: zypher encrypt [options] <input-value>
available options:
	-k, --key=<key>				key to encrypt/decrypt
	-f, --file=<path-to-file>		input file to be encrypted
	-o, --out=<path-to-file>		output file to be created
	-kf, --key-file=<path-to-file>  	A path to key file. Default: zypher.key
`
	SynopsisMsg = "encrypts input value or file with the provided key and prints the encrypted value to stdout or create a file"
)

type EncryptCmd struct {
	base BaseCmd
}

func NewEncryptCmd(cf CipherFactory, opts ...func(*BaseCmd)) *EncryptCmd {
	fs := flag.NewFlagSet("encrypt", flag.ContinueOnError)
	cfg := &config.Config{}
	fs.StringVar(&cfg.Key, "key", "", "key to encrypt/decrypt")
	fs.StringVar(&cfg.Key, "k", "", "key to encrypt/decrypt (shorthand)")
	fs.StringVar(&cfg.KeyFile, "key-file", "zypher.key", "file path for reading key to be used for encryption/decryption")
	fs.StringVar(&cfg.KeyFile, "kf", "zypher.key", "file path for reading key to be used for encryption/decryption (shorthand)")
	fs.StringVar(&cfg.OutFile, "out", "", "output file to be created")
	fs.StringVar(&cfg.OutFile, "o", "", "output file to be created (shorthand)")
	fs.StringVar(&cfg.InputFile, "file", "", "input file to be encrypted")
	fs.StringVar(&cfg.InputFile, "f", "", "input file to be encrypted (shorthand)")

	e := &EncryptCmd{
		base: BaseCmd{
			cfg: cfg,
			fs:  fs,
			cf:  cf,
			frw: file.NewFileReaderWriter(),
		},
	}

	for _, opt := range opts {
		opt(&e.base)
	}

	return e
}

func (e *EncryptCmd) Help() string {
	return HelpMsg
}

func (e *EncryptCmd) Synopsis() string {
	return SynopsisMsg
}

func (e *EncryptCmd) Run(args []string) int {
	if err := e.base.init(args); err != nil {
		fmt.Printf("error initializing encrypt cmd: %v", err)
		return 1
	}

	var (
		input []byte
		err   error
	)

	if e.base.cfg.Input != "" {
		input = []byte(e.base.cfg.Input)
	}

	if e.base.cfg.InputFile != "" {
		input, err = e.base.frw.ReadFile(e.base.cfg.InputFile)
		if err != nil {
			fmt.Printf("error reading input file: %v\n", err)
			return 1
		}
	}

	encrypted, err := e.base.ci.Encrypt(input)
	if err != nil {
		fmt.Printf("error encrypting: %v\n", err)
		return 1
	}

	base64Encrypted := base64.StdEncoding.EncodeToString(encrypted)
	if e.base.cfg.OutFile != "" {
		if err := e.base.frw.WriteFile(e.base.cfg.OutFile, []byte(base64Encrypted), 0600); err != nil {
			fmt.Printf("error writing to output file: %v\n", err)
			return 1
		}
		return 0
	} else {
		fmt.Println(base64Encrypted)
	}

	return 0
}
