# Kubernetes cluster with pre-configured strongSwan CNI plugi

The script provisons an 3 nodes Kubernete cluster with VirtualBox
driver. *kubeadm* is used for bootstrap cluster.

## Quick startup

```
$ vagrant up
$ kubectl --kubeconfig admin.conf get nodes
NAME       STATUS    AGE       VERSION
master     Ready     3h        v1.7.2
minion-1   Ready     2h        v1.7.2
minion-2   Ready     3m        v1.7.2
# Make sure you see 3 nodes ready, it may takes a few minutes
# Or if a node not show up after awhile, you can try reload it
# example
# vagrant reload minion2

# Optinal to run kubectl proxy for UI dashboard on http://127.0.0.1:8001/ui
kubectl --kubeconfig admin.conf proxy

# Bring up 9 pods with a simple echo server which listen on port 5678
$ kubectl --kubeconfig admin.conf create -f echo-demo-pod.yaml
namespace "demo" created
deployment "demo" created

# Wait until all pods are up
$ kubectl --kubeconfig admin.conf get pods --namespace=demo
NAME                    READY     STATUS    RESTARTS   AGE
demo-1594119973-1nvrn   1/1       Running   0          1m
demo-1594119973-4fp4r   1/1       Running   0          1m
demo-1594119973-8j0kw   1/1       Running   0          1m
demo-1594119973-phd4f   1/1       Running   0          1m
demo-1594119973-t9jc8   1/1       Running   0          1m
demo-1594119973-v8t5t   1/1       Running   0          1m
demo-1594119973-x5087   1/1       Running   0          1m
demo-1594119973-x94x6   1/1       Running   0          1m
demo-1594119973-zgkcd   1/1       Running   0          1m

# Now, you can attach shell to pods to views it ip
$ vagrant ssh master
$ sudo docker ps | grep echo
$ sudo docker exec -it container-id /bin/sh
/ # ip addr
		1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
				link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
				inet 127.0.0.1/8 scope host lo
					 valid_lft forever preferred_lft forever
		3: eth0@if12: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
				link/ether 0a:58:ac:11:00:0b brd ff:ff:ff:ff:ff:ff
				inet 172.17.0.11/16 scope global eth0
					 valid_lft forever preferred_lft forever
				inet 10.173.0.6/32 scope global eth0
					 valid_lft forever preferred_lft forever
# With above result, `10.173.0.28` is the ip address that is reachable from any pod
# on any server through IPSec tunnel
# From the pod, we can hit any other pods' ip
# Let's try to hit pod 10.173.0.7
/ # curl 10.173.0.7:5678
Pong. You hit me, my hostname is demo-1594119973-8j0kw, from 10.173.0.8:45050
/ # curl 10.173.0.9:5678
Pong. You hit me, my hostname is demo-1594119973-phd4f, from 10.173.0.8:52928/
```

## Demo Video

[https://youtu.be/fwk6gODMGn0](https://youtu.be/fwk6gODMGn0)

# Detail

## Cluster detail

Nodes name and their address respectively are:

* master: 10.9.0.2
* minion1: 10.9.0.3
* minion2: 10.9.0.4


If those conflict with other existing VM you have on the system, please
temporarily should down the other VMs for trying thi demonstration.

Master node is untained so they can run pod as well.

## Version

This is tested on those versiosn. It's recomened to have those
up-to-date for this demo to be succesfully.

* Kubenetes: 1.7.2
* Vagrant: 1.8.7

# How to use

Bring up cluster with:

```
vagrant up
```

After awhile, the cluster should be ready. If they fail or throw error
such as timeout try to reload them such as:

```
vagrant reload
# Or reload a particular node
vagrant reload minion2
```

We can SSH into any of them with:

```
vagrant ssh master
vagrant ssh minion1
vagrant ssh minion2
```

We also expose `admin.conf` so you can control the cluster easily by
passing `--kubeconfig admin.conf` such as:

```
kubectl --kubeconfig admin.conf get nodes
```

Note that it may take few minute for cluster to come up, depend on how
fast the machine since we run 3 VM at a same time.

Note that the cluster using *Kubernetes 1.7.2* so you may want to make
sure your `kubectl` is at same version.

Or you can also ssh into master, under *ubuntu* user and use `kubectl`
on there. It's preconfigured, you don't have to specify `kubeconfig`.

```
vagrant ssh master
$ whoami
ubuntu
$ kubectl get nodes
```

We also pre-configured dashboard. Hence you can simply run this from
your host machine to access dashboard:

```
kubectl --kubeconfig admin.conf proxy
```

Then open browser at [http://127.0.0.1:8001/ui](http://127.0.0.1:8001/ui)

## Spin up pods:

From this vagrant directory, deploy 9 pods with this:

```
kubectl --kubeconfig admin.conf apply -f echo-demo-pod.yaml
```

This creates namespace `demo`, run 9 pods of a simple echo server which echo
the target hostname and visitor ip address.

Once all pods are ready, you can attach a shell into them, and pod
should be able to communicate freely with other pods use the virtual ip
from VPN. That's the IP address of subnet *10.173.0.0/16*.

## Bonus

From master, where strongSwan is run, you can check connection status:

```
$ sudo ipsec statusall
# It should show something like this

Virtual IP pools (size/online/offline):
  10.173.0.0/16: 65534/11/27
Listening IP addresses:
  10.0.2.15
  10.9.0.2
Connections:
      server:  %any...%any  IKEv2
      server:   local:  [server] uses pre-shared key authentication
      server:   remote: uses pre-shared key authentication
      server:   child:  172.17.0.0/16 10.173.0.0/16 === dynamic TUNNEL
Security Associations (11 up, 0 connecting):
      server[40]: ESTABLISHED 16 minutes ago, 10.9.0.2[server]...172.17.0.20[13060]
      server[40]: IKEv2 SPIs: 8d5c14607ae3dadd_i fb65329255f0c38f_r*, pre-shared key reauthentication in 37 minutes
      server[40]: IKE proposal: AES_CBC_128/HMAC_SHA2_256_128/PRF_HMAC_SHA2_256/CURVE_25519
      server{56}:  INSTALLED, TUNNEL, reqid 40, ESP SPIs: c4c35a8d_i c9cb71d9_o
      server{56}:  AES_CBC_128/HMAC_SHA2_256_128, 0 bytes_i, 0 bytes_o, rekeying in 11 minutes
      server{56}:   10.173.0.0/16 172.17.0.0/16 === 10.173.0.38/32
      server[39]: ESTABLISHED 16 minutes ago, 10.9.0.2[server]...10.9.0.4[7978]
      server[39]: IKEv2 SPIs: b5eacb8758329e07_i 16649630993a59aa_r*, pre-shared key reauthentication in 39 minutes
      server[39]: IKE proposal: AES_CBC_128/HMAC_SHA2_256_128/PRF_HMAC_SHA2_256/CURVE_25519
      server{63}:  INSTALLED, TUNNEL, reqid 39, ESP in UDP SPIs: c86512bc_i c37dedb6_o
      server{63}:  AES_CBC_128/HMAC_SHA2_256_128, 0 bytes_i, 0 bytes_o, rekeying in 12 minutes
      server{63}:   10.173.0.0/16 172.17.0.0/16 === 10.173.0.37/32
      server[38]: ESTABLISHED 16 minutes ago, 10.9.0.2[server]...172.17.0.21[13065]
```
