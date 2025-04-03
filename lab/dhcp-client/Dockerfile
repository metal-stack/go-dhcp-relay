FROM debian:12

RUN apt update && apt install iproute2 isc-dhcp-client -y

CMD ["dhclient", "-v"]
