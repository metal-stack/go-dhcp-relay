package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/metal-stack/go-dhcp-relay/config"
	"github.com/metal-stack/go-dhcp-relay/server"
	"github.com/spf13/cobra"
)

var log *slog.Logger

var rootCmd = &cobra.Command{
	Use:   "go-dhcp-relay",
	Short: "A simple dhcp relay implementation",
	Run: func(cmd *cobra.Command, args []string) {
		iface, _ := cmd.Flags().GetString("interface")
		if iface == "" {
			log.Error("no listening interface was specified")
			os.Exit(1)
		}

		count, _ := cmd.Flags().GetUint8("count")
		servers, _ := cmd.Flags().GetStringArray("dhcp-servers")
		if len(servers) < 1 {
			log.Error("no dhcp servers were specified")
			os.Exit(1)
		}

		config := &config.Config{
			Interface:       iface,
			DHCPServers:     servers,
			MaximumHopCount: count,
		}
		if err := config.Validate(); err != nil {
			log.Error("invalid configuration", "error", err)
			os.Exit(1)
		}

		s, err := server.NewServer(log, config)
		if err != nil {
			log.Error("failed to initialize dhcp relay", "error", err)
			os.Exit(1)
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer func() {
			stop()
		}()

		var wg sync.WaitGroup
		var code int
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
			}()
			s.Serve(ctx)
		}()

		wg.Wait()
		os.Exit(code)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	rootCmd.Flags().StringP("interface", "i", "", "listening interface")
	rootCmd.Flags().Uint8P("count", "c", config.DefaultMaximumHopCount, "maximum hop count")
	rootCmd.Flags().StringArrayP("dhcp-servers", "s", nil, "list of dhcp servers to forward dhcp packets to")
}
