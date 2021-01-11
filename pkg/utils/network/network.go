package network

import (
	"io/ioutil"
	"net"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	cmd "nenvoy.com/pkg/utils/cmd"
)

type NetPlan struct {
	Network struct {
		Ethernets map[string]Ethernet `yaml:"ethernets,omitempty"`
		Version   int                 `yaml:"version"`
	} `yaml:"network"`
}

type Ethernet struct {
	Iface Iface `yaml:",inline"`
}

type Iface struct {
	Addresses []string `yaml:"addresses,omitempty"`
	DHCP4     *bool    `yaml:"dhcp4,omitempty"`
}

var ipNotFoundErr = errors.New("IP Address of interface not found")

// UpdateNetplanStatic - appends to neplan a static ip address adapter
func UpdateNetplanStatic(adapter string, ip string, prefix string) (err error) {
	file := "/etc/netplan/00-installer-config.yaml"

	//Read in the network config file
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	//Unmarshal the file into a struct
	n := new(NetPlan)
	err = yaml.Unmarshal(buf, n)
	if err != nil {
		return err
	}

	//Add the new interface into the struct
	newIface := Iface{
		[]string{ip + "/" + prefix},
		nil,
	}
	n.Network.Ethernets[adapter] = Ethernet{newIface}

	// Marshal the struct into yaml form
	out, err := yaml.Marshal(n)
	if err != nil {
		return err
	}

	//Write the file back to the netplan config file
	err = ioutil.WriteFile(file, out, 0644)
	if err != nil {
		return err
	}

	// Apply the new configuration
	_, stderr, err := cmd.Output("netplan", "apply")
	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	return nil

}

// EnableBridgedTraffic - Enables bridged trafic on the worker nodes
func EnableBridgedTraffic() (err error) {

	// Run the modprobe command
	_, stderr, err := cmd.Output("modprobe", "br_netfilter")
	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	// Write to the k8s config file for ip tables to allow traffic
	config := "net.bridge.bridge-nf-call-iptables = 1 \nnet.ipv4.ip_forward = 1 \nnet.bridge.bridge-nf-call-ip6tables = 1"
	err = ioutil.WriteFile("/etc/sysctl.d/99-kubernetes-cri.conf", []byte(config), 0644)
	if err != nil {
		return err
	}

	_, stderr, err = cmd.Output("sysctl", "--system")
	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	return nil
}

// HostFileAdd - Add entry to etc/hosts file
func HostFileAdd(hostIPAddr net.IP, hostname string) (err error) {

	newEntry := hostIPAddr.String() + " " + hostname + "\n"

	etcHosts, err := ioutil.ReadFile("/etc/hosts")

	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}

	output := append([]byte(newEntry), etcHosts...)

	err = ioutil.WriteFile("/etc/hosts", output, 0777)

	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}

	return nil
}

// UpdateHostsFile - Update /etc/hosts with peer nodes
func UpdateHostsFile(nodes map[string]string) (err error) {

	for nodeName, nodeIP := range nodes {
		err := HostFileAdd(net.ParseIP(nodeIP), nodeName)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetInterface - Returns the interface name matching an ip address
func GetInterface(ipAddr net.IP) (ifaceName string, err error) {

	// Get interfaces of localhost
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Find addresses matched by interface by name
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", errors.Wrap(err, "Could not collect addresses")
		}

		for _, addr := range addrs {
			ifaceAddr, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return "", errors.Wrap(err, "Could not parse ip")
			}

			if ifaceAddr.Equal(ipAddr) {
				return iface.Name, nil
			}
		}

	}

	return "", nil
}

// InterfaceExists - Checks if an interface exists
func InterfaceExists(ifaceName string) (exists bool, err error) {
	// Get interfaces of localhost
	ifaces, err := net.Interfaces()

	if err != nil {
		return false, err
	}

	// Find addresses matched by interface by name
	for _, iface := range ifaces {
		if iface.Name == ifaceName {
			return true, nil
		}
	}

	return false, nil
}

// PollIP - polls an interface to check if it has been assigned an IP
func PollIP(ifaceName string, ipType string, timeout int) (ip net.IP, err error) {

	for c := 0; c < timeout; c++ {
		ip, _, err = GetIPAddr(ifaceName, ipType)
		if err != nil && err != ipNotFoundErr {
			return ip, err
		} else if ip != nil {
			return ip, nil
		}

		time.Sleep(1 * time.Second)
	}

	return ip, errors.New("Could not get IP from interface: " + ifaceName)

}

// GetIPAddr - Returns the IP address assigned to the given network interface
// ipType: ipv4, ipv6
func GetIPAddr(ifaceName string, ipType string) (net.IP, *net.IPNet, error) {

	// Get interfaces of localhost
	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed on GetIPAddr (1):")
	}

	// Find addresses matched by interface by name
	for _, iface := range ifaces {
		if iface.Name == ifaceName {
			IPAddrs, err := iface.Addrs()

			if err != nil {
				return nil, nil, errors.Wrap(err, "failed on GetIPAddr (2):")
			}

			// Extract IP address
			for _, Addr := range IPAddrs {
				ipAddr, ipNet, err := net.ParseCIDR(Addr.String())

				if err != nil {
					return nil, nil, errors.Wrap(err, "failed on GetIPAddr (3):")
				}

				// Match by IP type
				if ipAddr.To4() != nil && ipType == "ipv4" {
					return ipAddr, ipNet, nil
				} else if ipAddr.To16() != nil && ipType == "ipv6" {
					return ipAddr, ipNet, nil
				}
			}
		}
	}

	return nil, nil, ipNotFoundErr
}

// Exists - Check if IP address is included in slice of IP addresses
func IPExists(ipAddr net.IP, ipAddrs []net.IP) bool {
	for _, ip := range ipAddrs {
		if ip.Equal(ipAddr) {
			return true
		}
	}
	return false
}
