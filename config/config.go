package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interface      string   `yaml:"interface"`
	GatewayAddress string   `yaml:"gateway-address"`
	DHCPServers    []string `yaml:"dhcp-servers"`
}

func UnmarshalConfig(in []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
