# StrongSwan CNI Plugin

A plugin to esbalish a pod-to-pod communication over IPSec, with virtual
ip return from strongSwan

# How it works

It's a modification of `bridge` plugin. We use `bridge` plugin, and
`ipam` to assign an ip address normally.

After network connections are up, we run `ipsec` inside network
namespace of container.

Every pods becomes a client of strongSwan, which is deployed separately,
and get an ip address from virtual ip pool of strongswan. The ip is then
assing to the `eth0` interface. All clients(pods in our case) can to each
others using that virtual ip.

## Components

### On Master

StrongSwan has to be pre-configured and run in master node. Currently,
the plugin only support IKEV2 PSK. We need to define a virtual ip pool,
pod will get ip address from this virtual ip pool.

## On Minion

Every pods from minion will connect to `strongSwan`, via IP address of
master server.

# How to use it

## Requirement on all nodes:

The host has to have `strongSwan` preinstalled so `ipsec` binary can be invoke.
StrongSwan can be install with this commands.

```
wget http://download.strongswan.org/strongswan-5.5.3.tar.bz2
tar xvf strongswan-5.5.3.tar.bz2
cd strongswan-5.5.3
sudo apt install build-essential libgmp-dev
/configure --prefix=/usr --sysconfdir=/etc --with-piddir=/etc/ipsec.d/run
make && sudo make install
```

Notice that here, we build strongswan outselves from source, because we want to
set a custom `piddir`. This custom `piddir` enable us to run multiple charon
instances.

## Requirement on master

On master, we need to run strongSwan as a daemon, it can run directly on host,
or as pod(in privileges mode, foward port 500 and 4500) up to developer devcison.
As long as we have a strongSwan server, we're fine.

## Install

* Puts the `strongswan` plugin executable file into `/opt/cni/bin/`. The file can be
download in relase page.

* Create a file `/etc/cni/net.d/10-swan.json` with this content

		```
		{
			"name": "ipsec",
			"type": "strongswan",
			"vpn": {
				// this is ip address where we run strongswan
				"serverIP": "10.9.0.2", 
				// this is virtual subnet that we will get ip address from. Same value on strongSwan server
				"virtualSubnet": "10.173.0.0/16",
				// this is subnet of minion.
				"hostSubnet": "10.9.0.0/24",
				// IKEV2 PSK: same value on strongSwan server
				"PSK": "dummy1234"
			},
			"bridge": "docker0",
			"isDefaultGateway": true,
			"forceAddress": false,
			"ipMasq": true,
			"hairpinMode": true,
			"ipam": {
				"type": "host-local",
				"subnet": "172.17.0.0/16"
			}
		}
		```

* Restart `kubelet` to take effective of this CNI plugin

# Demo

This is a demo video: To be added

# Trying out with Vagrant

We have a `Vagrantfile` which setup a sample cluster at: https://github.com/yeolabs/k8s-cni-ipsec/tree/master/final/vagrant
with this plugin. Follow instruction there to bring it up.
