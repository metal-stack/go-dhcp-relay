package cmd

import (
	"context"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		iface, _ := cmd.Flags().GetString("interface")
		if iface == "" {
			return fmt.Errorf("no listening interface was specified")
		}

		count, _ := cmd.Flags().GetUint8("count")
		servers, _ := cmd.Flags().GetStringArray("dhcp-servers")
		if len(servers) < 1 {
			return fmt.Errorf("no dhcp servers were specified")
		}

		config := &config.Config{
			Interface:       iface,
			DHCPServers:     servers,
			MaximumHopCount: count,
		}
		if err := config.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

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
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error("dhcp-relay", "error", err)
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
