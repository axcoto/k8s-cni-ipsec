#!/bin/bash -l

set -e

# Strongswan with custom piddir to run multiple one
wget http://download.strongswan.org/strongswan-5.5.3.tar.bz2
tar xvf strongswan-5.5.3.tar.bz2
cd strongswan-5.5.3
sudo apt install build-essential libgmp-dev
./configure --prefix=/usr --sysconfdir=/etc --with-piddir=/etc/ipsec.d/run
make && sudo make install


# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
apt-cache policy docker-ce
sudo apt-get install -y docker-ce
sudo usermod -aG docker ${USER}
sudo usermod -aG docker username
