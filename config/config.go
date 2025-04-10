package config

import (
	"fmt"
	"net"
)

const (
	DefaultMaximumHopCount = 4
)

type Config struct {
	Interface       string   `yaml:"interface"`
	DHCPServers     []string `yaml:"dhcp-servers"`
	MaximumHopCount uint8    `yaml:"maximum-hop-count"`
}

func (c *Config) Validate() error {
	if c.MaximumHopCount < 1 || c.MaximumHopCount > 16 {
		return fmt.Errorf("maximum hop count must be in range [1,16]")
	}

	iface, err := net.InterfaceByName(c.Interface)
	if err != nil {
		return fmt.Errorf("failed to retrieve information for interface %s:%w", c.Interface, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return fmt.Errorf("failed to retrieve ip addresses for interface %s:%w", c.Interface, err)
	}

	var ip4 net.IP
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return fmt.Errorf("failed to parse cidr %s:%w", addr.String(), err)
		}

		if ip.To4() != nil {
			ip4 = ip
		}
	}

	if ip4 == nil {
		return fmt.Errorf("no ipv4 address configured for interface %s", c.Interface)
	}

	for _, serverIP := range c.DHCPServers {
		if ip := net.ParseIP(serverIP); ip.To4() == nil {
			return fmt.Errorf("dhcp server address %s is not a valid ipv4 address", serverIP)
		}
	}

	return nil
}
