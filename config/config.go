package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	defaultMaximumHopCount = 4
)

type Config struct {
	Interface       string   `yaml:"interface"`
	GatewayAddress  string   `yaml:"gateway-address"`
	DHCPServers     []string `yaml:"dhcp-servers"`
	MaximumHopCount uint8    `yaml:"maximum-hop-count"`
}

func UnmarshalConfig(in []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SetDefaults() {
	if c.MaximumHopCount == 0 {
		c.MaximumHopCount = defaultMaximumHopCount
	}
}

func (c *Config) validate() error {
	if c.MaximumHopCount > 16 {
		return fmt.Errorf("maximum hop count must be in range [1,16]")
	}

	return nil
}
