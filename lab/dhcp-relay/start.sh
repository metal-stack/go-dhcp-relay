#!/bin/sh

apk add tcpdump
tcpdump -n -i any -w /etc/go-dhcp-relay/dump.pcap port 67 or port 68&
/usr/bin/go-dhcp-relay
