package config

type Config struct {
	Key       string
	KeyFile   string
	OutFile   string
	Input     string
	InputFile string
}

type ServerConfig struct {
	Port           int
	RootPubKeyPath string
}
