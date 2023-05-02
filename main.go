package zypher

import (
	"flag"
	"os"

	"github.com/mitchellh/cli"
	"github.com/vtno/zypher/cmd"
)

func main() {
	c := cli.NewCLI("zypher", "0.0.1")
	// TODO: read "key" from file > env > flag
	key := flag.String("key", "", "key to encrypt/decrypt")

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"encrypt": func() (cli.Command, error) {
			return cmd.NewEncryptCmd(NewCipher(*key)), nil
		},
	}
}
