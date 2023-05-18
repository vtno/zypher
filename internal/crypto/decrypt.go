package crypto

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/vtno/zypher/internal/config"
	"github.com/vtno/zypher/internal/file"
)

type DecryptCmd struct {
	base BaseCmd
}

const (
	DecryptHelpMsg = `Usage: zypher decrypt [options] <input-value>
available options:
	-k, --key=<key>			key for encryption/decryption
	-f, --file=<path-to-file>	input file to be encrypted
	-o, --out=<path-to-file>	output file to be created
`
	DecryptSynopsis = "decrypts input value or file with the provided key and prints the decrypted value to stdout or a file"
)

func NewDecryptCmd(cf CipherFactory, opts ...func(*BaseCmd)) *DecryptCmd {
	fs := flag.NewFlagSet("decrypt", flag.ContinueOnError)
	cfg := &config.Config{}
	fs.StringVar(&cfg.Key, "key", "", "key for encryption/decryption")
	fs.StringVar(&cfg.Key, "k", "", "key for encryption/decryption (shorthand)")
	fs.StringVar(&cfg.KeyFile, "key-file", "zypher.key", "file path for reading key to be used for encryption/decryption")
	fs.StringVar(&cfg.KeyFile, "kf", "zypher.key", "file path for reading key to be used for encryption/decryption (shorthand)")
	fs.StringVar(&cfg.OutFile, "out", "", "output file to be created")
	fs.StringVar(&cfg.OutFile, "o", "", "output file to be created (shorthand)")
	fs.StringVar(&cfg.InputFile, "file", "", "input file to be encrypted")
	fs.StringVar(&cfg.InputFile, "f", "", "input file to be encrypted (shorthand)")
	d := &DecryptCmd{
		base: BaseCmd{
			cfg: cfg,
			fs:  fs,
			cf:  cf,
			frw: file.NewFileReaderWriter(),
		},
	}
	for _, opt := range opts {
		opt(&d.base)
	}
	return d
}

func (d *DecryptCmd) Help() string {
	return DecryptHelpMsg
}

func (d *DecryptCmd) Synopsis() string {
	return DecryptSynopsis
}

func (d *DecryptCmd) Run(args []string) int {
	if err := d.base.init(args); err != nil {
		fmt.Printf("error initializing decrypt cmd %v", err)
		return 1
	}

	var (
		input string
		err   error
	)
	if d.base.cfg.Input != "" {
		input = d.base.cfg.Input
	}

	if d.base.cfg.InputFile != "" {
		fileInput, err := d.base.frw.ReadFile(d.base.cfg.InputFile)
		if err != nil {
			fmt.Printf("error reading input file: %v\n", err)
			return 1
		}
		input = string(fileInput)
	}

	decodedInput, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		fmt.Printf("error decoding input: %v\n", err)
		return 1
	}
	decrypted, err := d.base.ci.Decrypt(decodedInput)
	if err != nil {
		fmt.Printf("error decrypting: %v\n", err)
		return 1
	}

	if d.base.cfg.OutFile != "" {
		if err := d.base.frw.WriteFile(d.base.cfg.OutFile, decrypted, 0600); err != nil {
			fmt.Printf("error writing to file: %v\n", err)
			return 1
		}
	} else {
		fmt.Println(string(decrypted))
	}

	return 0
}
