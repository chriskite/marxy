#!/bin/bash

#
# haproxy.sh
#

HAPROXY="/etc/haproxy"
PIDFILE="/var/run/haproxy.pid"

cd "$HAPROXY"
haproxy -f /etc/haproxy/haproxy.cfg -p "$PIDFILE" 2>&1; while(true); do sleep 30; done
