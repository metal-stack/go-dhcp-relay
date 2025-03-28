package cmd

import (
	"log/slog"
	"os"

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

		err = server.Serve(log, config)
		if err != nil {
			log.Error("failed to start dhcp relay", "error", err)
			os.Exit(1)
		}
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
