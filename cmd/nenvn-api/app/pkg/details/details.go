package details

import (
	"encoding/json"

	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"
)

func GetHosts() (resp []byte, err error) {

	// Get all the hosts from the database
	hosts, err := host.GetHosts()
	if err != nil {
		return nil, err
	}

	// put the information into a JSON file
	resp, err = json.Marshal(hosts)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

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
