#!/bin/sh

sleep 5
dhclient -v

while :
do
  sleep 5
  dhclient -r -v
  dhclient -v
done
