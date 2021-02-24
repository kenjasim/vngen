package details

import (
	"encoding/json"

	"nenvoy.com/pkg/deployment"

	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"
)

// HostDetails - the details of a host to be viewed or presented
type HostDetails struct {
	Name       string
	Image      string
	State      string
	RAM        int
	CPUs       int
	Username   string
	Password   string
	HDSpace    string
	Deployment string
}

// GetHosts - Return all host details
func GetHosts() (resp []byte, err error) {

	// Get all the hosts from the database
	hosts, err := host.GetHosts()
	if err != nil {
		return nil, err
	}

	data := []HostDetails{}

	for _, host := range hosts {
		// Get the host details
		details, err := getHostDetails(host)
		if err != nil {
			return nil, err
		}

		data = append(data, details)
	}

	// put the information into a JSON file
	resp, err = json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetHost - Return a host's details
func GetHost(name string) (resp []byte, err error) {

	// Get all the hosts from the database
	host, err := host.GetHostByName(name)
	if err != nil {
		return nil, err
	}

	// Get the host details
	data, err := getHostDetails(host)

	// put the information into a JSON file
	resp, err = json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetHostIP - Returns the IP of a host
func GetHostIP(name string) (resp []byte, err error) {

	// Get all the hosts from the database
	host, err := host.GetHostByName(name)
	if err != nil {
		return nil, err
	}

	// Get the IP of the host
	ifaces, err := host.GetHostIfaces()

	// put the information into a JSON file
	resp, err = json.Marshal(ifaces)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetNetworks - Returns the Network of a host
func GetNetworks() (resp []byte, err error) {

	// Get all the networks from the database
	networks, err := network.GetNetworks()
	if err != nil {
		return nil, err
	}

	// put the information into a JSON file
	resp, err = json.Marshal(networks)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getHostDetails(host host.Host) (hostDet HostDetails, err error) {
	// Get state
	state, err := host.GetHostState()
	if err != nil {
		return hostDet, err
	}

	// Get deployment name
	dep, err := deployment.GetDeploymentByID(host.DeploymentID)
	if err != nil {
		return hostDet, err
	}

	// Assign to struct
	hostDet = HostDetails{Name: host.Name,
		Image:      host.Image,
		State:      state,
		RAM:        host.RAM,
		CPUs:       host.CPUs,
		Username:   host.Username,
		Password:   host.Password,
		HDSpace:    host.HDSpace,
		Deployment: dep.Name,
	}

	return hostDet, nil

}
