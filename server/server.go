package server

import (
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

func Serve(log *slog.Logger, config *config.Config) error {
	laddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: dhcpv4.ServerPort,
	}

	log.Info("starting dhcp-relay", "listen address", laddr)

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
	}()

	// TODO: add ctx and check if it's done
	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		packet, err := dhcpv4.FromBytes(buf)
		if err != nil {
			return err
		}

		log.Debug("received request", "read bytes", n, "peer", remoteAddr, "packet", packet.Summary())
	}
}
