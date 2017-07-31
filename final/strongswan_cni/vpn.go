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

// Establish an IPSec connection with strongSwan so that we can get an virtual IP
// This is a hack currently. Basically the interface wasn't ready when CNI run yet
// Hence we basically run a 10 times retry to bring up ipsec
// TODO: Rewrite this to avoid depend on binary ipsec and ip tool on the host
// We need a way to establish ipsec connection manually with strongswan
// Maybe need to look into libstrongswan
func establishIpsec(netNs string, containerId string, vpnInfo vpnInfo) error {
	netNs = extractProcId(netNs)
	log.Println("strongswan", "establish ipsec for", netNs)

	// Prepare directory tree
	// We're using ip netns, which require the network namespace in /var/run/netns/namespace
	// docker doesn't do this neither K8S, so we manually extract proc id and create symbol link
	os.Mkdir("/var/run/netns", os.ModePerm)
	os.Symlink(fmt.Sprintf("/proc/%s/ns/net", netNs), fmt.Sprintf("/var/run/netns/ns-%s", netNs))

	// When charon run, it puts pid file in /etc/ipsec.d/run hence we cannot run multiple instance
	// Luckily it has a capability to bind mount anything in /etc/netns/namespace/ into /etc/
	// respectively. We use this trick to create directory hold those pid and socket file
	os.Mkdir("/etc/ipsec.d/run", os.ModePerm)
	os.MkdirAll("/etc/netns/ns-"+netNs+"/ipsec.d/run", os.ModePerm)

	// Finally, generate client VPN configuration
	if err := genVpnConfig(netNs, vpnInfo); err != nil {
		return err
	}

	// Everything is ready, we can officially bring up ipsec
	args := []string{"bash", "-c", fmt.Sprintf(bringupIpsecScript, netNs), "&", "&>/tmp/nohup.log"}
	cmd := exec.Command("nohup", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("strongswan", "ipsec result", out.String())
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

// Extract procid to and use its as namespace in symlink
//  Example: /proc/27273/ns/net/ -> 27273
func extractProcId(netNs string) string {
	part := strings.Split(netNs, "/")
	return part[2]
}

// When CNI runs, the interface wasn't configured and up yet, we sleep a bit and re-try ten time before give up
const bringupIpsecScript = "for r in {1..10}; do sleep 10; if ip netns exec ns-%s ip addr | grep eth0; then ip netns exec ns-%s ipsec start >/dev/null 2>&1; break; fi; done"
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
