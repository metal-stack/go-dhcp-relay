#!/bin/bash

ip route add 10.255.0.3 via 10.1.255.2
dhcpd -f -d eth0
