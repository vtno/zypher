package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/vtno/zypher"
	"github.com/vtno/zypher/internal/crypto"
)

func main() {
	c := cli.NewCLI("zypher", "0.0.1")
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("zypher")
	c.Commands = map[string]cli.CommandFactory{
		"encrypt": func() (cli.Command, error) {
			return crypto.NewEncryptCmd(zypher.NewCipherFactory()), nil
		},
		"decrypt": func() (cli.Command, error) {
			return crypto.NewDecryptCmd(zypher.NewCipherFactory()), nil
		},
	}
	c.Run()
}
