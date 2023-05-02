package cmd

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

const (
	HelpMsg = `Usage: zypher encrypt [options] <input-value>
available options:
	-k, --key=<key>				key to encrypt/decrypt
	-f, --file=<path-to-file>	input file to be encrypted
	-o, --out=<path-to-file>	output file to be created
`
	SynopsisMsg = "encrypts input value or file with the provided key and prints the encrypted value to stdout or create a file"
)

type EncryptCmd struct {
	ci Cipher
}

func NewEncryptCmd(ci Cipher) *EncryptCmd {
	return &EncryptCmd{
		ci: ci,
	}
}

func (e EncryptCmd) Help() string {
	return HelpMsg
}

func (e EncryptCmd) Synopsis() string {
	return SynopsisMsg
}

func (e EncryptCmd) Run(args []string) int {
	return 0
}
