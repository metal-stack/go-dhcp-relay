package server

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (s *Server) handlePacket(packet *dhcpv4.DHCPv4) error {
	messageType := packet.MessageType()
	switch messageType {
	case dhcpv4.MessageTypeDiscover:
		return s.handleMessageTypeDiscover(packet)
	case dhcpv4.MessageTypeOffer:
		return s.handleMessageTypeOffer(packet)
	case dhcpv4.MessageTypeRequest:
		return s.handleMessageTypeRequest(packet)
	case dhcpv4.MessageTypeAck:
		return s.handleMessageTypeAck(packet)
	}

	// TODO: handle all types of messages

	return nil
}

func (s *Server) handleMessageTypeDiscover(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle discover", "packet", packet.Summary())
	if packet.HopCount >= 3 { // TODO: make this value configurable, rfc maximum is 16
		return fmt.Errorf("maximum hop count exceeded, dropping packet")
	}
	packet.HopCount++

	if packet.GatewayIPAddr.Equal(net.IPv4(0, 0, 0, 0)) {
		// TODO: check if correct: rfc "relay agent MUST fill this field with the IP address of the interface on which the request was received"
		packet.GatewayIPAddr = net.ParseIP(s.config.GatewayAddress)
	}
	packet.SetBroadcast()

	for _, serverIP := range s.config.DHCPServers {
		addr := &net.UDPAddr{
			IP:   net.ParseIP(serverIP),
			Port: dhcpv4.ServerPort,
		}

		n, err := s.sendTo(packet, addr)
		if err != nil {
			// FIX: returning here without trying the second is probably bad
			return fmt.Errorf("failed to send packet to server %s:%w", serverIP, err)
		}
		s.log.Debug("packet sent to server", "bytes sent", n, "server address", serverIP, "packet", packet.Summary())
	}

	return nil
}

func (s *Server) handleMessageTypeOffer(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle offer", "packet", packet.Summary())

	broadcastAddresses, err := getBroadcastAddresses(s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to get broadcast addresses for interface:%w", err)
	}

	for _, ip := range broadcastAddresses {
		addr := &net.UDPAddr{
			IP:   ip,
			Port: dhcpv4.ClientPort,
		}

		n, err := s.sendTo(packet, addr)
		if err != nil {
			return fmt.Errorf("failed to send packet to %s:%w", addr, err)
		}
		s.log.Debug("packet sent to client", "bytes sent", n, "address", addr, "packet", packet.Summary())
	}

	return nil
}

func (s *Server) handleMessageTypeRequest(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle request", "packet", packet.Summary())

	packet.GatewayIPAddr = net.ParseIP(s.config.GatewayAddress)

	serverIP := packet.ServerIdentifier()
	addr := &net.UDPAddr{
		IP:   serverIP,
		Port: dhcpv4.ServerPort,
	}

	n, err := s.sendTo(packet, addr)
	if err != nil {
		return fmt.Errorf("failed to send packet to server %s:%w", serverIP, err)
	}
	s.log.Debug("packet sent to server", "bytes sent", n, "server address", serverIP, "packet", packet.Summary())

	return nil
}

func (s *Server) handleMessageTypeAck(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle acknowledgment", "packet", packet.Summary())

	broadcastAddresses, err := getBroadcastAddresses(s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to get broadcast addresses for interface:%w", err)
	}

	for _, ip := range broadcastAddresses {
		addr := &net.UDPAddr{
			IP:   ip,
			Port: dhcpv4.ClientPort,
		}

		n, err := s.sendTo(packet, addr)
		if err != nil {
			return fmt.Errorf("failed to send packet to %s:%w", addr, err)
		}
		s.log.Debug("packet sent to client", "bytes sent", n, "address", addr, "packet", packet.Summary())
	}

	return nil
}

func (s *Server) sendTo(packet *dhcpv4.DHCPv4, addr *net.UDPAddr) (int, error) {
	n, err := s.conn.WriteToUDP(packet.ToBytes(), addr)
	if err != nil {
		return 0, err
	}

	return n, nil
}
