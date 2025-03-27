package cmd

import (
	"log/slog"
	"net"
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/spf13/cobra"
)

var log *slog.Logger

func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Debug("handle", "packet", *m, "peer", peer)
	client := client4.NewClient()

	conversation, _ := client.Exchange("wlan0", dhcpv4.WithServerIP(net.ParseIP("172.31.250.4")))
	for _, packet := range conversation {
		log.Debug("dhcp packet received", "packet", packet.Summary())
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-dhcp-relay",
	Short: "A simply dhcp relay implementation",
	Run: func(cmd *cobra.Command, args []string) {
		laddr := &net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: dhcpv4.ServerPort,
		}
		server, err := server4.NewServer("", laddr, handler)
		if err != nil {
			log.Error("unable to start server", "error", err)
			os.Exit(1)
		}

		log.Info("starting dhcp relay", "listen address", laddr)
		server.Serve()
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
}
