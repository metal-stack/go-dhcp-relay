# go-dhcp-relay

A simple dhcp relay implementation

## Usage

Run `go-dhcp-relay -c /path/to/go-dhcp-relay/config.yaml`.
If you don't specify a config path, `go-dhcp-relay` will expect the file to be located under `/etc/go-dhcp-relay/config.yaml`.

## Configuration

The following parameters can be configured.

| Field           | Type     | Description                    |
| --------------- | -------- | ------------------------------ |
| interface       | string   | the interface to configure     |
| gateway-address | string   | the relay's ip address         |
| dhcp-servers    | []string | the dhcp servers' ip addresses |
