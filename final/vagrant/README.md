# Kubernetes cluster with pre-configured strongSwan CNI plugi

The script provisons an 3 nodes Kubernete cluster with VirtualBox
driver. *kubeadm* is used for bootstrap cluster.

## Quick startup

Bring whole cluster with our default configuration, then deploy
9 pods of our echo server. Find ip address of pod, note them,
attach shell to a node and ping the above ip address.

```
make up
# Until all pods is in running state
make get_pods
make find_ip
make shell
/# curl http://10.173.0.2:5678
/# curl http://10.173.0.3:5678
/# curl http://10.173.0.4:5678
/# curl http://10.173.0.5:5678
/# curl http://10.173.0.6:5678
/# curl http://10.173.0.7:5678
/# curl http://10.173.0.8:5678
/# curl http://10.173.0.9:5678
```

# How this cluster is setup

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
passing `--kubeconfig admin.conf` direcly from your host machine:

```
kubectl --kubeconfig admin.conf get nodes
```

Note that the cluster using *Kubernetes 1.7.2* so you may want to make
sure your `kubectl` is at same version.

Or you can also ssh into master and use `kubectl` on there.

```
vagrant ssh master
$ sudo kubectl --kubeconfig /etc/kubernetes/admin.conf get nodes
```

This create namespace `demo` and deploy echo server under that.

## Spin up pods:

From this vagrant directory, deploy 9 pods with this:

```
kubectl --kubeconfig admin.conf apply -f echo-demo-pod.yaml
```

This creates namespace `demo`, run 9 pods of a simple echo server which echo
the target hostname and visitor ip address. This echo server runs on port 5678.
The source code of this server is at: https://github.com/yeolabs/k8s-cni-ipsec/tree/master/final/echo_server

Once all pods are ready, you can attach a shell into them, and try to
`curl` other pod ip.

They should be able to communicate freely with other pods use the virtual ip
from VPN. That's the IP address of subnet *10.173.0.0/16*.

To find the IP Address of node, you can SSH into a server and find the docker
container then use `docker exec` to run `ip addr`

```
vagrant ssh master
sudo docker ps | grep demo
sudo docker exec -it [above-container-id] ip addr
```

Once you know IP address of all pod, simply attach shell to a pod and request other pods:

```
vagrant ssh master
sudo docker exec -it [container-id] shell
$ curl [other-pod-ip]:5678
```

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


## Demo Video

[https://youtu.be/fwk6gODMGn0](https://youtu.be/fwk6gODMGn0)
