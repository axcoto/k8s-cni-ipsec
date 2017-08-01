#!/usr/bin/env bash

set -xe

#=============================================================#
# Try to prepopluation some env if not seting. Vagrant script
# will set those env
PROVISION_DIR="/vagrant/provison"
KUBEADM_TOKEN="${KUBEADM_TOKEN:-390699.8859a9a052d8b9f9}"
MASTER_ADDR="${MASTER_ADDR:-10.9.0.2}"
NODE_ROLE="$1"
#=============================================================#

# Install latest docker to prepare for K8S using it
install_docker() {
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  sudo apt-get update
  apt-cache policy docker-ce
  sudo apt-get install -y docker-ce
  sudo usermod -aG docker ${USER}
}

# Install k8s stack with kubeadm. We also add kubectl binary on every node
install_kube() {
  curl -sLO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
  chmod +x ./kubectl
  sudo mv ./kubectl /usr/local/bin/kubectl

  sudo apt-get install -y apt-transport-https
  curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
  cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
  deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF

  sudo apt-get update
  sudo apt-get install -y kubelet kubeadm
}

# Build a custom strongSwan instead of using default ont from APT
#
# We need to customize `piddir` into `/etc/ipsec.d` so that network namespace
# can bind mount thing in /etc/netns/<namespace> 
# Reference: https://wiki.strongswan.org/projects/strongswan/wiki/Netns
install_strongswan() {
  wget http://download.strongswan.org/strongswan-5.5.3.tar.bz2
  tar xvf strongswan-5.5.3.tar.bz2
  cd strongswan-5.5.3
  sudo apt install -y build-essential libgmp-dev pkg-config libsystemd-dev
  ./configure --prefix=/usr --sysconfdir=/etc --with-piddir=/etc/ipsec.d/run --enable-systemd --enable-swanctl
  make && sudo make install
  cd ..
  rm -rf strongswan-5.5.3 strongswan-5.5.3.tar.bz2
  mkdir -p etc/ipsec.d/run
  echo '%any : PSK dummy1234' > /etc/ipsec.secrets
}

# Generate ipsec config.
# We only need to run this on master
config_strongswan() {
  sudo cp "$PROVISION_DIR/ipsec.conf" /etc/ipsec.conf
  sudo systemctl enable strongswan
  sudo systemctl start strongswan
}

# Install,config strongswan cni plugin
config_cni() {
  mkdir -p /etc/cni/net.d
  cat <<EOF >/etc/cni/net.d/10-strongswan.json
{
	"name": "ipsec",
	"type": "strongswan",
	"vpn": {
		"serverIP": "$MASTER_ADDR",
		"virtualSubnet": "10.173.0.0/16",
		"hostSubnet": "10.9.0.0/24",
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
EOF
	sudo curl -sL https://github.com/yeolabs/k8s-cni-ipsec/releases/download/0.1-rc1/strongswan -o /opt/cni/bin/strongswan
	sudo chmod +x /opt/cni/bin/strongswan
	sudo systemctl restart kubelet
}

# Using kubeadm to initialize cluster and prepare necessary config
# to work with our cluster
create_cluster() {
  sudo kubeadm init  --apiserver-advertise-address="$MASTER_ADDR"
  sudo kubeadm token create "$KUBEADM_TOKEN" --ttl 0
  config_cni
  cp /etc/kubernetes/admin.conf /vagrant/

	mkdir -p $HOME/.kube
	sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
	sudo chown $(id -u):$(id -g) $HOME/.kube/config
	kubectl taint nodes --all node-role.kubernetes.io/master-
}

# Using kubeadm to join existing cluster and prepare necessay config for cni
join_cluster() {
  sudo kubeadm join --token "$KUBEADM_TOKEN" "$MASTER_ADDR":6443 --skip-preflight-checks
	config_cni
}

init_master() {
  config_strongswan
  create_cluster
}

init_secondary() {
  join_cluster
}

# main entry point.
main() {
  install_docker
  install_kube
  install_strongswan

  if [ "$NODE_ROLE" = "master" ]; then
  	init_master
  else
    init_secondary
  fi
}

# Fire it up!!!
main
