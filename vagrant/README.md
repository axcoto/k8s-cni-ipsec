# Kubernetes cluster with pre-configured strongSwan CNI plugi

The script provisons an 3 nodes Kubernete cluster with VirtualBox
driver. *kubeadm* is used for bootstrap cluster.

## Quick startup

```
vagrant up
kubectl --kubeconfig admin.conf get nodes
# Make sure you see 3 nodes

# Optinal to run kubectl proxy for UI dashboard
kubectl --kubeconfig admin.conf proxy
# Now, the UI should be accessbible a thttp://127.0.0.1:8001/ui

# Bring up 9 pods with a simple echo server which listen on port 5678
kubectl --kubeconfig admin.conf create -f echo-demo-pod.yaml
# Wait until all pods are up

kubectl --kubeconfig admin.conf get nodes --namespace=demo
# Attach shell to pods to views it ip
# You should be able to `curl` any of those ip on any pods
kubectl --kubeconfig admin.conf --namespace=demo exec -it demo -c 3295696929-23jvr -- sh
# Now you are inside that pods, view the ip address, do the samething
for other pods and try to curl them on port 5678
ip addr
curl 10.173.0.10:5678
```

## Demo Video



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
on there

```
vagrant ssh master
$ whoami
ubuntu
$ kubectl get nodes
```

We also pre-configured dashboard. Hence you can simply run this from
your host machine:

```
kubectl --kubeconfig admin.conf proxy
```

Then open browser at [http://127.0.0.1:8001/ui](http://127.0.0.1:8001/ui)

## Spin up pods:

From this vagrant directory, deploy 9 pods with this:

```
kubectl --kubeconfig apply -f echo.yaml
```

Once all pods are ready, you can attach a shell into them, and pod
should be able to communicate freely with other pods use the virtual ip
from VPN. That's the IP address of subnet *10.173.0.0/16*.
