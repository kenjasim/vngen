package network

import (
	"encoding/xml"

	libvirt "libvirt.org/libvirt-go"

	structs "nenvoy.com/pkg/constants"
)

// CreateNetworkXML - Create the network definition file
func CreateNetworkXML(net structs.NetworkDefinition) (domainDef string, err error) {
	//Define the domain object for libvirt
	network := structs.Network{}

	// Set the name of the network
	network.Name = net.NetworkName

	// Set the forward mode
	network.Forward.Mode = net.Type

	// Sort the bridge
	network.Bridge.Name = net.NetworkName
	network.Bridge.Stp = "on"
	network.Bridge.Delay = "0"

	// Ip Address setup
	network.IP.Address = net.NetworkAddr
	network.IP.Netmask = net.Netmask

	// Setup the DHCP server
	network.IP.Dhcp.Range.Start = net.DHCPLower
	network.IP.Dhcp.Range.End = net.DHCPUpper

	xmlBytes, err := xml.MarshalIndent(network, "", "	")
	if err != nil {
		return "", err
	}

	return string(xmlBytes), nil
}

// CreateNetwork - Apply the network xml config and create the network
func CreateNetwork(networkDef string) (err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create the network
	network, err := conn.NetworkDefineXML(networkDef)
	if err != nil {
		return err
	}

	// Set the network to autostart
	err = network.SetAutostart(true)
	if err != nil {
		return err
	}

	// Create the network
	err = network.Create()
	if err != nil {
		return err
	}

	return nil
}
