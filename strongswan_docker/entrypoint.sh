#!/bin/sh -e
#
# entrypoint for strongswan
#
# env |grep vpn_ | while read line; do echo $line| cut -d= -f2- >> /etc/ipsec.d/secrets.local.conf ; done

INTERFACE=${IPTABLES_INTERFACE:+-i ${IPTABLES_INTERFACE}}
ENDPOINTS=${IPTABLES_ENDPOINTS:+-s ${IPTABLES_ENDPOINTS}}
RIGHTSUBNET=$(grep rightsubnet /etc/ipsec.docker/ipsec.gc.conf  | cut -d"=" -f2)

# enable ip forward
sysctl -w net.ipv4.ip_forward=1

# function to use when this script recieves a SIGTERM.
_term() {
  echo "Caught SIGTERM signal! Stopping ipsec..."
  #kill -TERM "$child" 2>/dev/null
  ipsec stop
}

# catch the SIGTERM
trap _term SIGTERM

echo "Starting strongSwan/ipsec..."
ipsec start --nofork "$@" &

child=$!
# wait for child process to exit
wait "$child"
