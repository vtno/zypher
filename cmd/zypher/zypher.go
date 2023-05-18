package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/vtno/zypher"
	"github.com/vtno/zypher/internal/crypto"
	"github.com/vtno/zypher/internal/keygen"
)

const version = "0.2.0"

func main() {
	c := cli.NewCLI("zypher", version)
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("zypher")
	c.Commands = map[string]cli.CommandFactory{
		"encrypt": func() (cli.Command, error) {
			return crypto.NewEncryptCmd(zypher.NewCipherFactory()), nil
		},
		"decrypt": func() (cli.Command, error) {
			return crypto.NewDecryptCmd(zypher.NewCipherFactory()), nil
		},
		"keygen": func() (cli.Command, error) {
			return keygen.NewKeyGenCmd(), nil
		},
	}
	_, err := c.Run()
	if err != nil {
		log.Fatalf("error running zypher command: %s", err)
	}
}
