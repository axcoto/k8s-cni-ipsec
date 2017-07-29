# StrongSwan CNI Plugin

A plugin to esbalish a pod-to-pod communication over IPSec, with virtual
ip return from strongswan

# How it works

It's an modification of `bridge` plugin. We use `bridge` plugin, to
create interface.

After network connections are up, we run `ipsec` inside network
namespace of container.

Every pods become a client of strongswan, and get an ip address from
virtual ip pool of strongswan.

## Components

### On Master

On master, we run a pod in `strongswan` namespace, in privilegesd mode,
forward port 400 and 4500 into pods

## On Minion

Every pods from minion will connect to `strongswan`, via IP address of
master server.

# How to use it

## Requirement on all nodes:

The host has to have `strongswan` preinstalled so `ipsec` binary can be invoke.
StrongSwan can be install with this commands.

```
wget http://download.strongswan.org/strongswan-5.5.3.tar.bz2
tar xvf strongswan-5.5.3.tar.bz2
cd strongswan-5.5.3
sudo apt install build-essential libgmp-dev
/configure --prefix=/usr --sysconfdir=/etc --with-piddir=/etc/ipsec.d/run
make && sudo make install
```

Notice that here, we build strongswan outselve from source, because we want to set
a custom `piddir`. This custom `piddir` enable us to run multiple charon instances.

## Requirement on master

On master, the strongswan daemon has to be run, it can run directly on host,
or as pod(in privileges mode, foward port 500 and 4500) up to developer devcison.
As long as we have a strongSwan server, we're fine

## Install

### Puts the `strongswan` plugin executable file into `/opt/cni/bin/`,
create a file `/etc/cni/net.d/10-swan.json` with this content

```
{
	"name": "mynet",
	"type": "bridge",
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

## Restart `kubelet` to take effective of this CNI plugin


# Demo

This is a demo video:

# Trying out with Vagrant

We have a `Vagrantfile` which expose cluster at: https://github.com/yeolabs/k8s-cni-ipsec/tree/master/vagrant

Clone this repository, cd into vagrant folder and run:

```
vagrant up
```

This provision a cluster of 3 nodes, on private subnet `10.9.0.0/24`.

Once cluster is up, you can ssh into any node:

```
vagrant ssh master
vagrant ssh minion1
vagrant ssh minion2
```

Now, let's deploy something. Inside `vagrant` filder, run this:

```
kubectl --kubeconfig admin.conf apply nginx-spec.yaml
```

Or you can use `dashboard`:

```
kubectl --kubeconfig admin.conf proxy
```

And open browser at : http://127.0.0.1:8001/ui

After pods are ready. You can attach a shell into nodes, and check their ip address:

```
ip addr
```

You can freely ping any ip between them or run `curl [the-private-ip]`
