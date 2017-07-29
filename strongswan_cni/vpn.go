package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	ipsecSecretPath = "/etc/ipsec.secrets"
)

// TODO: Rewrite this to avoid depend on binary ipsec and ip tool on the host
// We need a way to establish ipsec connection manually with strongswan
// Maybe need to look into libstrongswan
func establishIpsec(netNs string, containerId string, vpnInfo vpnInfo) error {
	netNs = extractProcId(netNs)
	log.Println("strongswan", "establish ipsec for", netNs)

	os.Mkdir("/var/run/netns", os.ModePerm)
	os.Mkdir("/etc/ipsec.d/run", os.ModePerm)
	// Directory to hold charon pid file, this will be bindmount to /etc/ipsec.d/run in netowkr namespace
	os.MkdirAll("/etc/netns/ns-"+netNs+"/ipsec.d/run", os.ModePerm)

	os.Symlink(fmt.Sprintf("/proc/%s/ns/net", netNs), fmt.Sprintf("/var/run/netns/ns-%s", netNs))
	// Create ipsec.conf file
	//cp /etc/ipsec.client /etc/netns/ns-$pid/ipsec.conf
	//sed -i s/@leftid/@container-pid-$pid/ /etc/netns/ns-$pid/ipsec.conf
	// We use netNamesapce as leftid so that if container get kill, it gets namespace
	// of pause pods and will get same virtual ip
	if err := genVpnConfig(netNs, vpnInfo); err != nil {
		return err
	}

	ipinfo, _ := exec.Command("ip", "netns", "exec", "/sbin/ip", "addr").Output()
	log.Println("strongswan", "ipaddr", string(ipinfo[:]))

	// Bringup ipsec
	args := []string{"bash", "-c", fmt.Sprintf("sleep 20; ip netns exec ns-%s ipsec start >>/tmp/cni-swan.log 2>&1", netNs), "&", "&>/tmp/nohup.log"}
	cmd := exec.Command("nohup", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("strongswan", "ipsec result", out.String())
	//ip netns exec ns-$pid ipsec start
	return nil
}

// Stop ipsec, clearout namespace/configfile,symbol link that we have set
func teardownIpsec(netNs string) {
	netNs = extractProcId(netNs)
	log.Println("strongswan", "teardown ipsec for", netNs)
	exec.Command("ip", "netns", "exec", "ns-"+netNs, "ipsec", "stop").Run()
}

// Generate VPN config for pod
func genVpnConfig(netNs string, vpnInfo vpnInfo) error {
	configContent := ipsecConf
	configContent = strings.Replace(configContent, "$LeftId$", "@"+netNs, 1)
	configContent = strings.Replace(configContent, "$ServerIP$", vpnInfo.ServerIP, 1)
	configContent = strings.Replace(configContent, "$VirtualSubnet$", vpnInfo.VirtualSubnet, 1)
	configContent = strings.Replace(configContent, "$HostSubnet$", vpnInfo.HostSubnet, 1)

	os.MkdirAll("/etc/netns/ns-"+netNs, os.ModePerm)
	if err := ioutil.WriteFile("/etc/netns/ns-"+netNs+"/ipsec.conf", []byte(configContent), 0644); err != nil {
		return err
	}

	if _, err := os.Stat(ipsecSecretPath); os.IsNotExist(err) {
		ioutil.WriteFile(ipsecSecretPath, []byte(fmt.Sprintf("%%any : PSK %s", vpnInfo.PSK)), 0644)
	}

	return nil
}

const ipsecConf = `conn %default
	ikelifetime=60m
	keylife=20m
	rekeymargin=3m
	keyingtries=1
	keyexchange=ikev2
	authby=secret

conn home
	left=%any
	leftsourceip=%config
	leftid=$LeftId$
	leftfirewall=yes
	right=$ServerIP$
	rightsubnet=172.17.0.0/16,$VirtualSubnet$,$HostSubnet$
	rightid=server
	auto=start`
