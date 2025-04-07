package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/metal-stack/go-dhcp-relay/config"
)

type Server struct {
	config *config.Config
	log    *slog.Logger
	conn   *net.UDPConn
}

func NewServer(log *slog.Logger, config *config.Config) (*Server, error) {
	laddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: dhcpv4.ServerPort,
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		config: config,
		log:    log,
		conn:   conn,
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

	n, err := s.conn.Read(bytes)
	if err != nil {
		errChan <- fmt.Errorf("failed to read message: %w", err)
		return
	}
	s.log.Debug("message received", "bytes read", n)
	recvChan <- bytes
}
