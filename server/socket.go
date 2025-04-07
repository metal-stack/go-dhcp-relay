package server

import (
	"encoding/binary"
	"net"
	"net/netip"
)

func getBroadcastAddresses(ifname string) ([]net.IP, error) {
	intf, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	addrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, 0)
	for _, addr := range addrs {
		ip, err := broadcastFromCIDR(addr)
		if err != nil {
			return nil, err
		}
		if ip != nil {
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

func broadcastFromCIDR(addr net.Addr) (net.IP, error) {
	ip, cidr, err := net.ParseCIDR(addr.String())
	if err != nil {
		return nil, err
	}
	if ip.To4() == nil {
		return nil, nil
	}
	cidrInt := binary.BigEndian.Uint32(cidr.Mask)
	suffix := cidrInt ^ 0xffffffff // flip the bits

	a, err := netip.ParseAddr(cidr.IP.String())
	if err != nil {
		return nil, err
	}

	ipInt := binary.BigEndian.Uint32(a.AsSlice())
	broadcastInt := ipInt | suffix // all ones where mask was 0

	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, broadcastInt)

	return net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3]), nil
}
