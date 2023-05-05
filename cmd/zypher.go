package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/vtno/zypher"
	"github.com/vtno/zypher/internal/encrypt"
)

type CipherFactory struct{}

func (cf *CipherFactory) NewCipher(key string) encrypt.Cipher {
	return zypher.NewCipher(key)
}

var (
	key       string
	outFile   string
	inputFile string
)

func main() {
	c := cli.NewCLI("zypher", "0.0.1")
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("zypher")
	c.Commands = map[string]cli.CommandFactory{
		"encrypt": func() (cli.Command, error) {
			return encrypt.NewEncryptCmd(&CipherFactory{}), nil
		},
	}
	c.Run()
}
