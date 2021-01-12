package network

import (
	"encoding/xml"
	"fmt"

	"nenvoy.com/pkg/utils/printing"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	libvirt "libvirt.org/libvirt-go"

	"nenvoy.com/pkg/constants"
	structs "nenvoy.com/pkg/constants"
)

//Network - Struct for the network data in the database
type Network struct {
	gorm.Model
	Name         string
	IP           string
	DHCPLower    string
	DHCPUpper    string
	Netmask      string
	Type         string
	DeploymentID uint
}

// createNetworkXML - Create the network definition file
func (n *Network) createNetworkXML() (domainDef string, err error) {
	//Define the domain object for libvirt
	network := structs.Network{}

	// Set the name of the network
	network.Name = n.Name

	// Set the forward mode
	network.Forward.Mode = n.Type

	// Sort the bridge
	network.Bridge.Name = n.Name
	network.Bridge.Stp = "on"
	network.Bridge.Delay = "0"

	// Ip Address setup
	network.IP.Address = n.IP
	network.IP.Netmask = n.Netmask

	// Setup the DHCP server
	network.IP.Dhcp.Range.Start = n.DHCPLower
	network.IP.Dhcp.Range.End = n.DHCPUpper

	xmlBytes, err := xml.MarshalIndent(network, "", "	")
	if err != nil {
		return "", err
	}

	return string(xmlBytes), nil
}

// CreateNetwork - Apply the network xml config and create the network
func (n *Network) CreateNetwork() (err error) {

	networkDef, err := n.createNetworkXML()
	if err != nil {
		return err
	}
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

// Destroy - Destroy the network
func (n *Network) Destroy() (err error) {

	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	network, err := conn.LookupNetworkByName(n.Name)
	if err != nil {
		return err
	}

	// Destroy the network
	err = network.Destroy()
	if err != nil {
		return err
	}

	// Undefine the network
	err = network.Undefine()
	if err != nil {
		return err
	}

	printing.PrintSuccess(fmt.Sprintf("Destroyed network %s", n.Name))

	return nil
}

// DefineNetwork - Defines the network struct to be added to the database and creates the xml file
func DefineNetwork(net structs.NetworkDefinition) (network Network, err error) {
	// Create network struct for database
	network = Network{
		Name:      net.NetworkName,
		IP:        net.NetworkAddr,
		DHCPLower: net.DHCPLower,
		DHCPUpper: net.DHCPUpper,
		Netmask:   net.Netmask,
		Type:      net.Type,
	}

	if err != nil {
		return network, err
	}

	return network, nil
}

// GetNetworks - returns all the networks in the database
func GetNetworks() (networks []Network, err error) {
	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return networks, errors.Wrap(err, "failed to connect database")
	}

	err = db.Find(&networks).Error
	if err != nil {
		return networks, errors.Wrap(err, "could not find networks")
	}

	return networks, nil
}

//GetNetworksByDeployment - returns all the networks in a deployment
func GetNetworksByDeployment(ID uint) (networks []Network, err error) {
	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return networks, errors.Wrap(err, "failed to connect database")
	}

	err = db.Where("deployment_id = ?", ID).Find(&networks).Error
	if err != nil {
		return networks, errors.Wrap(err, "could not find networks")
	}

	return networks, nil
}
