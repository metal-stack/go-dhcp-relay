FROM debian:12

RUN apt update && apt install iproute2 isc-dhcp-server -y

RUN touch /var/lib/dhcp/dhcpd.leases
COPY dhcpd.conf /etc/dhcp/dhcpd.conf

ENTRYPOINT ["dhcpd", "-f", "eth0"]
