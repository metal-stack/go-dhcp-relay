FROM debian:12

RUN apt update && apt install iproute2 isc-dhcp-server tcpdump iputils-ping -y

RUN touch /var/lib/dhcp/dhcpd.leases
COPY dhcpd.conf /etc/dhcp/dhcpd.conf
COPY start.sh /usr/bin/start.sh
RUN chmod +x /usr/bin/start.sh

ENTRYPOINT ["/usr/bin/start.sh"]
