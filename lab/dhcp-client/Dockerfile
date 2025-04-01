FROM debian:12

RUN apt update && apt install isc-dhcp-client -y

ENTRYPOINT ["dhclient", "-v"]
