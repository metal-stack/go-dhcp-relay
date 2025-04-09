#!/bin/sh

sleep 5

ip addr add 10.1.255.18/24 dev eno2
ip route add 10.255.0.0/24 dev eno2
dhcpd -f -d eno2
