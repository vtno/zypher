package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/vtno/zypher/internal/config"
	"github.com/vtno/zypher/internal/server/auth"
	"github.com/vtno/zypher/internal/server/provider"
	"github.com/vtno/zypher/internal/server/store"
	"go.uber.org/zap"
)

type ServerCmd struct{}

const (
	HelpMsg = `Usage: zypher server [options]
    -p, --port          a port to start the server. default: 8080
		    --rootKeyPath		a path of root public key. default: ~/.ssh/id_rsa.pub
    `
	Synopsis = "starts a key server"
)

func NewServerCmd() *ServerCmd {
	return &ServerCmd{}
}

func (s *ServerCmd) Help() string {
	return HelpMsg
}

func (s *ServerCmd) Synopsis() string {
	return Synopsis
}

func (s *ServerCmd) Run(arg []string) int {
	cfg := &config.ServerConfig{}
	ctx := context.Background()

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.IntVar(&cfg.Port, "port", 8080, "a port to start the server")
	fs.StringVar(&cfg.RootPubKeyPath, "rootKeyPath", "zypher.pub", "a path of root public key")

	bbStore, err := store.NewBBoltStore(defaultDbPath)
	if err != nil {
		fmt.Printf("error creating bbStore: %v", err)
		return 1
	}

	pkProvider := provider.NewPubKeyProvider(cfg.RootPubKeyPath)
	a, err := auth.NewAuth(bbStore, auth.WithPubKeyProvider(pkProvider))
	if err != nil {
		fmt.Printf("error creating auth: %v", err)
		return 1
	}

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("error creating a logger: %v", err)
		return 1
	}

	srv, err := NewServer(bbStore, a, logger)
	if err != nil {
		fmt.Printf("error creating a server %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Println("stopping zypher server")
		srv.Stop(ctx)
	}()

	fmt.Printf("starting zypher server at port %d\n", cfg.Port)
	if err := srv.Start(); err != http.ErrServerClosed {
		fmt.Printf("error starting a server %v", err)
	}

	return 0
}
