package server

import (
	"context"
	"errors"
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

type Server struct {
	config *config.Config
	log    *slog.Logger
}

func NewServer(log *slog.Logger, config *config.Config) *Server {
	return &Server{
		config: config,
		log:    log,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	laddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: dhcpv4.ServerPort,
	}

	s.log.Info("starting dhcp-relay", "listen address", laddr)

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
	}()

	packet := make(chan *dhcpv4.DHCPv4)

	for {
		go func() {
			p, err := s.listenForPackets(conn)
			if err != nil && !errors.Is(err, net.ErrClosed) {
				s.log.Error("error occurred while listening for packets", "error", err)
				return
			}

			packet <- p
		}()

		select {
		case <-ctx.Done():
			s.log.Info("shutdown signal received")
			return nil

		case p := <-packet:
			err := s.handlePacket(p)
			if err != nil {
				s.log.Error("error occurred while processing packet", "error", err)
			}
		}
	}
}

func (s *Server) listenForPackets(conn *net.UDPConn) (*dhcpv4.DHCPv4, error) {
	buf := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFrom(buf)
	if err != nil {
		return nil, err
	}

	p, err := dhcpv4.FromBytes(buf)
	if err != nil {
		return nil, err
	}

	s.log.Debug("packet received", "read bytes", n, "peer", remoteAddr, "packet", p.Summary())
	return p, nil
}

func (s *Server) handlePacket(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle packet", "packet", packet.Summary())
	return nil
}
