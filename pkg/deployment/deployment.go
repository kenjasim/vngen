package deployment

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	constants "nenvoy.com/pkg/constants"
	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"
	"nenvoy.com/pkg/utils/printing"
)

//Deployment - Struct for the deployment data in the database
type Deployment struct {
	gorm.Model
	ID       uint
	Name     string
	Hosts    []host.Host
	Networks []network.Network
}

// CreateDeployment - Create the deployment and relevent hosts
func CreateDeployment(vnDef constants.VirtualNetworkDefinition) (err error) {

	printing.PrintInfo(fmt.Sprintf("Creating deployment %s...", vnDef.Deployment.DeploymentName))

	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}

	// Migrate the database if not already
	err = db.AutoMigrate(&host.Host{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	err = db.AutoMigrate(&network.Network{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	err = db.AutoMigrate(&Deployment{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	//Create the deployment from the virtual network definition
	dep := &Deployment{Name: vnDef.Deployment.DeploymentName}

	// Create the networks
	err = CreateNetworks(vnDef, dep)
	if err != nil {
		return errors.Wrap(err, "failed to create networks: ")
	}

	// Create the hosts
	err = CreateHosts(vnDef, dep)
	if err != nil {
		return errors.Wrap(err, "failed to create hosts: ")
	}

	// Write the deployment to the database
	db.Create(dep)

	return nil
}

// CreateNetworks - Creates the networks in KVM and adds them to the deployment
func CreateNetworks(vnDef constants.VirtualNetworkDefinition, dep *Deployment) (err error) {
	printing.PrintInfo("Creating networks...")

	// Create the network definitions and create them
	for _, netwr := range vnDef.Networks {
		net, networkDef, err := network.DefineNetwork(netwr)
		if err != nil {
			return err
		}

		// Create the network
		err = network.CreateNetwork(networkDef)
		if err != nil {
			return err
		}

		// Append the network to the deployment
		dep.Networks = append(dep.Networks, net)

		printing.PrintSuccess(fmt.Sprintf("Created %s network %s with ip %s", netwr.Type, netwr.NetworkName, netwr.NetworkAddr))
	}

	return nil
}

// CreateHosts - Creates the hosts in KVM and adds them to the deployment
func CreateHosts(vnDef constants.VirtualNetworkDefinition, dep *Deployment) (err error) {

	printing.PrintInfo("Creating hosts...")

	// Create the host definition files and create the hosts
	for _, hst := range vnDef.Host {
		hostDB, hostDef, err := host.DefineHost(hst)
		if err != nil {
			return err
		}

		err = host.CreateHost(hostDef)
		if err != nil {
			return err
		}

		// Append the host to the deployment
		dep.Hosts = append(dep.Hosts, hostDB)

		printing.PrintSuccess(fmt.Sprintf("Created host %s", hst.HostName))
	}

	return nil
}
