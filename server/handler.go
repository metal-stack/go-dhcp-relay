package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"golang.org/x/net/ipv4"
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
	case dhcpv4.MessageTypeDecline:
		return s.handleMessageTypeDecline(packet)
	case dhcpv4.MessageTypeInform:
		return s.handleMessageTypeInform(packet)
	case dhcpv4.MessageTypeNak:
		return s.handleMessageTypeNak(packet)
	}

	return nil
}

func (s *Server) handleMessageTypeDiscover(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle discover", "packet", packet.Summary())

	if packet.HopCount >= s.config.MaximumHopCount {
		return fmt.Errorf("maximum hop count exceeded, dropping packet")
	}
	packet.HopCount++

	return s.sendToAllServers(packet)
}

func (s *Server) handleMessageTypeInform(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle inform", "packet", packet.Summary())

	if packet.HopCount >= s.config.MaximumHopCount {
		return fmt.Errorf("maximum hop count exceeded, dropping packet")
	}
	packet.HopCount++

	return s.sendToAllServers(packet)
}

func (s *Server) handleMessageTypeRequest(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle request", "packet", packet.Summary())
	return s.sendToServer(packet)
}

func (s *Server) handleMessageTypeDecline(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle decline", "packet", packet.Summary())
	return s.sendToServer(packet)
}

func (s *Server) handleMessageTypeOffer(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle offer", "packet", packet.Summary())
	return s.sendBroadcastToClient(packet)
}

func (s *Server) handleMessageTypeNak(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle negative acknowledgment", "packet", packet.Summary())
	return s.sendBroadcastToClient(packet)
}

func (s *Server) handleMessageTypeAck(packet *dhcpv4.DHCPv4) error {
	s.log.Debug("handle acknowledgment", "packet", packet.Summary())

	if packet.ClientIPAddr.Equal(net.IPv4zero) {
		return s.sendBroadcastToClient(packet)
	}

	addr := &net.UDPAddr{
		IP:   packet.ClientIPAddr,
		Port: dhcpv4.ClientPort,
	}

	n, err := s.sendTo(packet, addr, &s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to send packet to %s:%w", addr, err)
	}
	s.log.Debug("packet sent to client", "bytes sent", n, "address", addr, "packet", packet.Summary())

	return nil
}

func (s *Server) sendBroadcastToClient(packet *dhcpv4.DHCPv4) error {
	addr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: dhcpv4.ClientPort,
	}

	n, err := s.sendTo(packet, addr, &s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to send packet to %s:%w", addr, err)
	}
	s.log.Debug("packet sent to client", "bytes sent", n, "address", addr, "packet", packet.Summary())

	return nil
}

func (s *Server) sendToServer(packet *dhcpv4.DHCPv4) error {
	err := injectGatewayAddress(packet, s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to inject gateway address:%w", err)
	}

	serverIP := packet.ServerIdentifier()
	addr := &net.UDPAddr{
		IP:   serverIP,
		Port: dhcpv4.ServerPort,
	}

	n, err := s.sendTo(packet, addr, nil)
	if err != nil {
		return fmt.Errorf("failed to send packet to server %s:%w", serverIP, err)
	}
	s.log.Debug("packet sent to server", "bytes sent", n, "server address", serverIP, "packet", packet.Summary())

	return nil
}

func (s *Server) sendToAllServers(packet *dhcpv4.DHCPv4) error {
	err := injectGatewayAddress(packet, s.config.Interface)
	if err != nil {
		return fmt.Errorf("failed to inject gateway address:%w", err)
	}

	errs := make([]error, 0)
	for _, serverIP := range s.config.DHCPServers {
		addr := &net.UDPAddr{
			IP:   net.ParseIP(serverIP),
			Port: dhcpv4.ServerPort,
		}

		n, err := s.sendTo(packet, addr, nil)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s:%w", serverIP, err))
		}
		s.log.Debug("packet sent to server", "bytes sent", n, "server address", serverIP, "packet", packet.Summary())
	}

	return errors.Join(errs...)
}

func (s *Server) sendTo(packet *dhcpv4.DHCPv4, addr *net.UDPAddr, ifname *string) (int, error) {
	var cm *ipv4.ControlMessage
	if ifname != nil {
		intf, err := net.InterfaceByName(*ifname)
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve interface:%w", err)
		}
		cm = &ipv4.ControlMessage{
			IfIndex: intf.Index,
		}
	}

	n, err := s.conn.WriteTo(packet.ToBytes(), cm, addr)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func injectGatewayAddress(packet *dhcpv4.DHCPv4, ifname string) error {
	if !packet.GatewayIPAddr.Equal(net.IPv4zero) {
		return nil
	}

	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return err
		}
		if ip.To4() != nil {
			packet.GatewayIPAddr = ip
		}
	}

	return nil
}
