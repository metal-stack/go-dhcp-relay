package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/metal-stack/go-dhcp-relay/config"
	"github.com/metal-stack/go-dhcp-relay/server"
	"github.com/urfave/cli/v2"
)

var (
	interfaceFlag = &cli.StringFlag{
		Name:    "interface",
		Aliases: []string{"i"},
		Value:   "",
		Usage:   "listening interface",
		EnvVars: []string{"INTERFACE"},
	}
	dhcpServersFlag = &cli.StringSliceFlag{
		Name:    "dhcp-servers",
		Aliases: []string{"s"},
		Value:   cli.NewStringSlice(),
		Usage:   "list of dhcp servers to forward requests to",
		EnvVars: []string{"DHCP_SERVERS"},
	}
	maximumHopCountFlag = &cli.UintFlag{
		Name:    "maximum-hop-count",
		Aliases: []string{"c"},
		Value:   config.DefaultMaximumHopCount,
		Usage:   "packets whose hop count exceeds this value will be dropped",
		EnvVars: []string{"MAXIMUM_HOP_COUNT"},
	}
)

func rootCmd(cCtx *cli.Context) error {
	iface := cCtx.String("interface")
	if iface == "" {
		return fmt.Errorf("no listening interface was specified")
	}

	count := cCtx.Uint("maximum-hop-count")
	servers := cCtx.StringSlice("dhcp-servers")
	if len(servers) < 1 {
		return fmt.Errorf("no dhcp servers were specified")
	}

	config := &config.Config{
		Interface:       iface,
		DHCPServers:     servers,
		MaximumHopCount: uint8(count),
	}
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	s, err := server.NewServer(log, config)
	if err != nil {
		return fmt.Errorf("failed to initialize dhcp relay:%w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		stop()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		s.Serve(ctx)
	}()

	wg.Wait()
	return nil
}

func main() {

	app := &cli.App{
		Name:   "go-dhcp-server",
		Usage:  "A simple dhcp relay implementation",
		Action: rootCmd,
		Flags: []cli.Flag{
			interfaceFlag,
			dhcpServersFlag,
			maximumHopCountFlag,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to start dhcp server:%W", err)
	}
}
