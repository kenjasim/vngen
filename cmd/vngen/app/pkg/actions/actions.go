package actions

import (
	"encoding/json"
	"fmt"

	"nenvoy.com/pkg/topology"

	structs "nenvoy.com/pkg/constants"
)

// Build - takes a JSON template file and builds it
func Build(template []byte) (err error) {
	vnDef := structs.VirtualNetworkDefinition{}

	// Unmarhsal the JSON template
	err = json.Unmarshal(template, &vnDef)
	if err != nil {
		return err
	}

	fmt.Println(vnDef)

	// Build the template
	err = topology.Build(vnDef)
	if err != nil {
		return err
	}

	return nil

}

// Start - Starts either a deployment or host
func Start(name string, resource string) (err error) {
	// Check if you want to start the host or deployment
	if resource == "host" {
		err = topology.StartHost(name)
		if err != nil {
			return err
		}
	} else if resource == "deployment" {
		err = topology.StartDeployment(name)
		if err != nil {
			return err
		}
	}

	return nil

}

// Stop - Stops either a deployment or host
func Stop(name string, resource string) (err error) {
	// Check if you want to stop the host or deployment
	if resource == "host" {
		err = topology.StopHost(name)
		if err != nil {
			return err
		}
	} else if resource == "deployment" {
		err = topology.StopDeployment(name)
		if err != nil {
			return err
		}
	}

	return nil

}

// Restart - Restarts either a deployment or host
func Restart(name string, resource string) (err error) {
	// Check if you want to restart the host or deployment
	if resource == "host" {
		err = topology.RestartHost(name)
		if err != nil {
			return err
		}
	} else if resource == "deployment" {
		err = topology.RestartDeployment(name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Destroy - Restarts either a deployment or host
func Destroy(name string, resource string) (err error) {
	// Check if you want to restart the host or deployment
	if resource == "host" {
		err = topology.DestroyHost(name)
		if err != nil {
			return err
		}
	} else if resource == "deployment" {
		err = topology.DestroyDeployment(name)
		if err != nil {
			return err
		}
	}

	return nil
}
