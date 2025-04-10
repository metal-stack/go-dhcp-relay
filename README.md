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

Run `make lab-up` to deploy the following containerlab topology.

![topology](images/go-dhcp-relay.clab.drawio)

Then you can inspect the container logs to see if everything is working as expected.
The expected behavior is that all DHCP requests sent by `dhcp-client1` and `dhcp-client2` will be forwarded to both DHCP servers, while the requests sent by `dhcp-client3` are ignored because it is connected to a different interface than the one the DHCP relay is configured to listen on.
As soon as a client has acquired a lease it will sleep for 5 seconds and then start another exchange.
