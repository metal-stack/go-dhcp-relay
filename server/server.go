package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/metal-stack/go-dhcp-relay/config"
	"golang.org/x/net/ipv4"
)

type Server struct {
	config *config.Config
	log    *slog.Logger
	conn   *ipv4.PacketConn
}

func NewServer(log *slog.Logger, config *config.Config) (*Server, error) {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf("%s:%d", net.IPv4zero, dhcpv4.ServerPort))
	if err != nil {
		return nil, err
	}
	packetConn := ipv4.NewPacketConn(conn)

	s := &Server{
		config: config,
		log:    log,
		conn:   packetConn,
	}

	return s, nil
}

func (s *Server) Serve(ctx context.Context) {
	defer func() {
		err := s.conn.Close()
		if err != nil {
			s.log.Error("error while closing connection", "error", err)
		}
	}()

	recvChan := make(chan []byte)
	errChan := make(chan error)

	for {
		go s.listen(recvChan, errChan)

		select {
		case <-ctx.Done():
			s.log.Debug("shutdown signal received")
			return
		case recv := <-recvChan:
			packet, err := dhcpv4.FromBytes(recv)
			if err != nil {
				s.log.Error("failed to parse packet", "error", err)
				continue
			}

			err = s.handlePacket(packet)
			if err != nil {
				s.log.Error("failed to process packet", "error", err)
			}
		case err := <-errChan:
			s.log.Error("error listening for packets", "error", err)
		}
	}
}

func (s *Server) listen(recvChan chan<- []byte, errChan chan<- error) {
	bytes := make([]byte, 1024)

	n, _, src, err := s.conn.ReadFrom(bytes)
	if err != nil {
		errChan <- fmt.Errorf("failed to read message: %w", err)
		return
	}
	s.log.Debug("message received", "bytes read", n, "source address", src)
	recvChan <- bytes
}
