package server

import (
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

type handler struct {
	log    *slog.Logger
	config *config.Config
}

func (h *handler) handle(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	h.log.Debug("handle dhcp request", "packet", *m, "peer", peer)
	for _, server := range h.config.DHCPServers {
		go h.exchange(server)
	}
}

func (h *handler) exchange(serverIP string) {
	client := client4.NewClient()

	conversation, err := client.Exchange(h.config.Interface, dhcpv4.WithServerIP(net.ParseIP(serverIP)))
	for _, packet := range conversation {
		h.log.Debug("dhcp packet received", "packet", packet.Summary())
	}

	if err != nil {
		h.log.Error("error occurred during dhcp exchange", "error", err)
	}
}
