package handler

import (
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

type handler struct {
	log    *slog.Logger
	config *config.Config
}

func NewHandler(log *slog.Logger, config *config.Config) *handler {
	return &handler{
		log:    log,
		config: config,
	}
}

func (h *handler) Handle(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	h.log.Debug("handle packet", "packet", m.Summary())
}
