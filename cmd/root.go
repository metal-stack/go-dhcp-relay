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
		configFile, _ := cmd.Flags().GetString("config")

		configBytes, err := os.ReadFile(configFile)
		if err != nil {
			log.Error("failed to open config file", "error", err)
			os.Exit(1)
		}

		config, err := config.UnmarshalConfig(configBytes)
		if err != nil {
			log.Error("failed to parse config file", "error", err)
			os.Exit(1)
		}
		config.SetDefaults()

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

	rootCmd.Flags().StringP("config", "c", "/etc/go-dhcp-relay/config.yaml", "path to config file")
}
