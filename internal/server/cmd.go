package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/vtno/zypher/internal/auth"
	"github.com/vtno/zypher/internal/config"
	"github.com/vtno/zypher/internal/store"
)

type ServerCmd struct{}

const (
	HelpMsg = `Usage: zypher server [options]
    -p, --port    a port to start the server. default: 8080
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

	bbStore, err := store.NewBBoltStore(defaultDbPath)
	if err != nil {
		fmt.Printf("error creating bbStore: %v", err)
		return 1
	}

	srv, err := NewServer(bbStore, auth.NewAuthGuard(bbStore))
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
