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

# How to use it

## Requirement on host:

The host has to installed `strongswan` so `ipsec` binary can be invoke.
StrongSwan can be install with this commands.

```

```

## Puts the `strongswan` plugin executable file into `/opt/cni/bin/`,
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
