# go-dhcp-relay

A simple DHCP relay implementation.

## Usage

```bash
go-dhcp-relay -i interface -s server1 [... -s serverN]
```

## Test Locally

### Requirements

- [docker](https://docs.docker.com/)
- [docker compose](https://docs.docker.com/compose/)
- [containerlab](https://containerlab.dev/)

### Starting Lab

Run `make lab-up` to deploy containerlab topology with one relay, two DHCP servers and three clients.
The relay contains two VLANs, Vlan4000 and Vlan1000, and is configured to listen for DHCP requests on the Vlan4000.
Two of the clients are connected to Vlan4000 on the relay, the third one is connected to Vlan1000.
The expected behavior is that the first two clients' requests are forwarded by the the relay, while the third client's requests are dropped.

Inspect the container logs to observe what is happening.
As soon as a client has acquired a lease or if it failed to do so it will sleep for 5 seconds and start another exchange.
