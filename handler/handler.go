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
	h.log.Debug("handle packet", "connection", conn, "peer", peer, "packet", m.Summary())

	if m.OpCode != dhcpv4.OpcodeBootRequest {
		h.log.Error("invalid opcode received, only bootrequest accepted", "message", m.Summary())
		return
	}

	if m.HopCount >= 3 {
		h.log.Info("dropping packet because hop count exceeded 3", "hops", m.HopCount)
		return
	}
	m.HopCount++

	if m.GatewayIPAddr.Equal(net.IPv4(0, 0, 0, 0)) {
		m.GatewayIPAddr = net.ParseIP(h.config.GatewayAddress)
	}
	m.SetBroadcast()

	for _, s := range h.config.DHCPServers {
		err := h.exchange(conn, m, s)
		if err != nil {
			h.log.Error("error during dhcp exchange occurred", "server", s, "error", err)
		}
	}
}

func (h *handler) exchange(clientConn net.PacketConn, m *dhcpv4.DHCPv4, serverIP string) error {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP(h.config.GatewayAddress),
		Port: dhcpv4.ClientPort,
	}

	raddr := &net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: dhcpv4.ServerPort,
	}

	serverConn, err := net.DialUDP("udp", laddr, raddr) // error bind address in use
	if err != nil {
		return err
	}
	defer func() {
		err := serverConn.Close()
		if err != nil {
			h.log.Error("error occurred while closing connection", "error", err)
		}
	}()

	_, err = serverConn.Write(m.ToBytes())
	if err != nil {
		return err
	}

	received := make([]byte, 1024)
	_, err = serverConn.Read(received)

	packet, err := dhcpv4.FromBytes(received)
	if err != nil {
		return err
	}

	h.log.Debug("reply received", "packet", packet.Summary())

	return nil
}
