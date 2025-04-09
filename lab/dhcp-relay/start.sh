#!/bin/sh

sleep 5

ip link add Bridge type bridge
ip link set Bridge type bridge vlan_filtering 1
ip link set Bridge up

ip link set Ethernet0 master Bridge
ip link set Ethernet1 master Bridge

ip link add link Bridge name Vlan4000 type vlan id 4000
ip addr add 10.255.0.1/24 dev Vlan4000
ip link set Vlan4000 up

bridge vlan add vid 4000 dev Ethernet0 pvid untagged
bridge vlan add vid 4000 dev Ethernet1 pvid untagged
bridge vlan add vid 4000 dev Bridge self pvid untagged

ip route add 10.1.255.2 dev eth0
ip route add 10.1.255.18 dev eth1

/usr/bin/go-dhcp-relay
