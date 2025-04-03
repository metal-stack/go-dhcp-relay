package server

import (
	"context"
	"errors"
	"fmt"
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

	for _, ip := range s.config.DHCPServers {
		// FIX: this should be done concurrently
		err := s.exchange(packet, ip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) exchange(packet *dhcpv4.DHCPv4, serverIP string) error {
	if packet.HopCount >= 3 {
		return fmt.Errorf("maximum hop count exceeded, hops: %d", packet.HopCount)
	}
	packet.HopCount++

	if packet.GatewayIPAddr.Equal(net.IPv4(0, 0, 0, 0)) {
		packet.GatewayIPAddr = net.ParseIP(s.config.GatewayAddress)
	}
	packet.SetBroadcast()

	relayClientAddr := &net.UDPAddr{
		IP:   net.ParseIP(s.config.GatewayAddress),
		Port: dhcpv4.ClientPort,
	}
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: dhcpv4.ServerPort,
	}

	err := s.sendTo(packet, relayClientAddr, serverAddr)
	if err != nil {
		return err
	}

	// FIX: this address is in use because relay is already listening there. pass connection to the function and use that one
	relayServerAddr := &net.UDPAddr{
		IP:   net.ParseIP(s.config.GatewayAddress),
		Port: dhcpv4.ServerPort,
	}

	_, err = s.receive(relayServerAddr, serverAddr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendTo(packet *dhcpv4.DHCPv4, laddr, raddr *net.UDPAddr) error {
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			s.log.Error("error while closing connection", "error", err)
		}
	}()

	n, err := conn.Write(packet.ToBytes())
	if err != nil {
		return err
	}
	s.log.Debug("sent packet to peer", "bytes sent", n, "local address", laddr, "remote address", raddr)

	return nil
}

func (s *Server) receive(laddr, raddr *net.UDPAddr) (*dhcpv4.DHCPv4, error) {
	var recv []byte
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return nil, err
	}

	// FIX: add timeout
	n, err := conn.Read(recv)
	if err != nil {
		return nil, err
	}
	s.log.Debug("received reply from peer", "bytes read", n, "local address", laddr, "remote address", raddr)

	packet, err := dhcpv4.FromBytes(recv)
	if err != nil {
		return nil, err
	}
	s.log.Debug("reply received", "packet", packet.Summary())

	return packet, nil
}
