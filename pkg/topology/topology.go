package topology

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"nenvoy.com/pkg/constants"
	"nenvoy.com/pkg/database"
	"nenvoy.com/pkg/deployment"
	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"
	"nenvoy.com/pkg/utils/printing"

	structs "nenvoy.com/pkg/constants"
)

// BuildFromFile - Allows the building of a VN from a file
func BuildFromFile(filename string) (err error) {
	printing.PrintInfo(fmt.Sprintf("Building  deployment %s...", filename))

	// Read in the yaml config file
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	vnDef := &structs.VirtualNetworkDefinition{}
	err = yaml.Unmarshal(buf, vnDef)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}

	err = Build(*vnDef)
	if err != nil {
		return err
	}

	return nil

}

//Build - Build the virtual network
func Build(vnDef structs.VirtualNetworkDefinition) (err error) {

	printing.PrintInfo(fmt.Sprintf("Creating deployment %s...", vnDef.Deployment.DeploymentName))

	// Connect and open the database
	db, err := database.NewSession()
	if err != nil {
		return err
	}

	// Ensure the hosts, networks and deployments are migrated
	err = migrateDatabase(db)
	if err != nil {
		return err
	}

	//Create the deployment from the virtual network definition
	dep := &deployment.Deployment{Name: vnDef.Deployment.DeploymentName}

	// Create the networks
	err = createNetworks(vnDef, dep)
	if err != nil {
		creationError := errors.Wrap(err, "failed to create networks")
		cleanupDeployment(dep)
		return creationError
	}

	// Create the hosts
	err = createHosts(vnDef, dep)
	if err != nil {
		creationError := errors.Wrap(err, "failed to create hosts")
		cleanupDeployment(dep)
		return creationError
	}

	// Write the deployment to the database
	db.Create(dep)

	return nil

}

// StartDeployment - Starts the deployment by name
func StartDeployment(depName string) (err error) {
	// Get the deployment
	dep, err := deployment.GetDeploymentByName(depName)
	if err != nil {
		return err
	}

	// Get the hosts which have the same deployment ID
	hosts, err := host.GetHostsByDeployment(dep.ID)
	if err != nil {
		return err
	}

	for _, hst := range hosts {
		// Start the host
		err := hst.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

// StartHost - Starts the host by name
func StartHost(name string) (err error) {
	// Get the hosts which have the same deployment ID
	hst, err := host.GetHostByName(name)
	if err != nil {
		return err
	}

	//Start the host
	err = hst.Start()
	if err != nil {
		return err
	}

	return nil
}

// RestartDeployment - Restarts the deployment by name
func RestartDeployment(depName string) (err error) {
	// Get the deployment
	dep, err := deployment.GetDeploymentByName(depName)
	if err != nil {
		return err
	}

	// Get the hosts which have the same deployment ID
	hosts, err := host.GetHostsByDeployment(dep.ID)
	if err != nil {
		return err
	}

	for _, hst := range hosts {
		// Start the host
		err := hst.Restart()
		if err != nil {
			return err
		}
	}

	return nil
}

// RestartHost - Restarts the host by name
func RestartHost(name string) (err error) {
	// Get the hosts which have the same deployment ID
	hst, err := host.GetHostByName(name)
	if err != nil {
		return err
	}

	//Start the host
	err = hst.Restart()
	if err != nil {
		return err
	}

	return nil
}

// StopDeployment - Stops the deployment by name
func StopDeployment(depName string) (err error) {
	// Get the deployment
	dep, err := deployment.GetDeploymentByName(depName)
	if err != nil {
		return err
	}

	// Get the hosts which have the same deployment ID
	hosts, err := host.GetHostsByDeployment(dep.ID)
	if err != nil {
		return err
	}

	for _, hst := range hosts {
		// Start the host
		err := hst.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}

// StopHost - Stops the host by name
func StopHost(name string) (err error) {
	// Get the hosts which have the same deployment ID
	hst, err := host.GetHostByName(name)
	if err != nil {
		return err
	}

	//Start the host
	err = hst.Stop()
	if err != nil {
		return err
	}

	return nil
}

// DestroyDeployment - Stops the deployment by name
func DestroyDeployment(depName string) (err error) {
	// Get the deployment
	dep, err := deployment.GetDeploymentByName(depName)
	if err != nil {
		return err
	}

	// Get the hosts which have the same deployment ID
	hosts, err := host.GetHostsByDeployment(dep.ID)
	if err != nil {
		return err
	}

	// Destroy hosts
	for _, hst := range hosts {
		// Start the host
		err := hst.Destroy()
		if err != nil {
			return err
		}
	}

	// Get the hosts which have the same deployment ID
	networks, err := network.GetNetworksByDeployment(dep.ID)
	if err != nil {
		return err
	}

	// Destroy networks
	for _, netwk := range networks {
		err := netwk.Destroy()
		if err != nil {
			return err
		}
	}

	// Destroy the deployment
	err = dep.Destroy()
	if err != nil {
		return err
	}

	return nil
}

//DestroyHost - Destroys a single host
func DestroyHost(name string) (err error) {
	// Get the hosts which have the same deployment ID
	hst, err := host.GetHostByName(name)
	if err != nil {
		return err
	}

	//Start the host
	err = hst.Destroy()
	if err != nil {
		return err
	}

	return nil
}

// migrateDatabase - ensures that the migrations have been applied
func migrateDatabase(db *gorm.DB) (err error) {
	// Migrate the database if not already
	err = db.AutoMigrate(&host.Host{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	err = db.AutoMigrate(&network.Network{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	err = db.AutoMigrate(&deployment.Deployment{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate database: ")
	}

	return nil
}

// createNetworks - Creates the networks in KVM and adds them to the deployment
func createNetworks(vnDef constants.VirtualNetworkDefinition, dep *deployment.Deployment) (err error) {
	printing.PrintInfo("Creating networks...")

	// Create the network definitions and create them
	for _, netwr := range vnDef.Networks {
		net, err := network.DefineNetwork(netwr)
		if err != nil {
			return err
		}

		// Create the network
		err = net.CreateNetwork()
		if err != nil {
			return err
		}

		// Append the network to the deployment
		dep.Networks = append(dep.Networks, net)

		printing.PrintSuccess(fmt.Sprintf("Created %s network %s with ip %s", netwr.Type, netwr.NetworkName, netwr.NetworkAddr))
	}

	return nil
}

// createHosts - Creates the hosts in KVM and adds them to the deployment
func createHosts(vnDef constants.VirtualNetworkDefinition, dep *deployment.Deployment) (err error) {

	printing.PrintInfo("Creating hosts...")

	// Create the host definition files and create the hosts
	for _, hst := range vnDef.Host {
		hostDB, err := host.DefineHost(hst)
		if err != nil {
			return err
		}

		err = hostDB.CreateHost(hst.Networks)
		if err != nil {
			return err
		}

		// Append the host to the deployment
		dep.Hosts = append(dep.Hosts, hostDB)

		printing.PrintSuccess(fmt.Sprintf("Created host %s", hst.HostName))
	}

	return nil
}

func cleanupDeployment(dep *deployment.Deployment) {
	printing.PrintWarning("Build failed, starting cleanup...")
	printing.PrintInfo("Starting cleanup due to error...")
	// Undefine hosts and networks
	for _, netwrk := range dep.Networks {
		// Undefine the networks
		netwrk.Destroy()
	}

	for _, hst := range dep.Hosts {
		// Undefine the networks
		hst.Destroy()
	}

	printing.PrintSuccess("Cleanup finished")
}
