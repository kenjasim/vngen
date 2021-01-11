package config

import (
	"errors"
	"net"
	"strings"

	"github.com/spf13/viper"
	structs "nenvoy.com/pkg/constants"
	"nenvoy.com/pkg/utils/network"
	"nenvoy.com/pkg/utils/printing"
)

// ReadConfig - Reads the application config file
func ReadConfig(cfgFileP *string) (err error) {
	cfgFile := *cfgFileP
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "config.yaml" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	//Read in the config file
	err = viper.ReadInConfig()

	if err != nil {
		printing.PrintWarning(err.Error())
	}

	err = detectCurrentNode()
	if err != nil {
		return err
	}

	return nil
}

//detectCurrentNode - detects current node and adds them to viper
func detectCurrentNode() (err error) {
	var clusterNodes []structs.ClusterNodes
	viper.UnmarshalKey("clusterNodes", &clusterNodes)

	for _, node := range clusterNodes {
		iface, err := network.GetInterface(net.ParseIP(node.NodeIPAddr))
		if err != nil {
			return err
		}

		if iface != "" {
			viper.Set("nodeDef.nodeIPAddr", node.NodeIPAddr)
			viper.Set("nodeDef.nodeName", node.NodeName)
			viper.Set("nodeDef.role", node.Role)
			return nil
		}
	}

	return errors.New("Could not detect host from IP address")

}

//UnmarshalNodes - Unmarshals the nodes into maps
func UnmarshalNodes(filter string) (nodes map[string]string) {

	nodes = make(map[string]string)

	var clusterNodes []structs.ClusterNodes
	viper.UnmarshalKey("clusterNodes", &clusterNodes)

	for _, node := range clusterNodes {
		if strings.ToLower(filter) == "all" || strings.ToLower(node.Role) == strings.ToLower(filter) {
			nodes[node.NodeName] = node.NodeIPAddr
		}
	}

	return nodes
}

//ExtractNodeIPs - Collect the ip addresses from the node definitions depending on a filter (all, master, worker)
func ExtractNodeIPs(filter string) []net.IP {
	var clusterNodes []structs.ClusterNodes
	viper.UnmarshalKey("clusterNodes", &clusterNodes)

	nodeIPs := []net.IP{}

	for _, node := range clusterNodes {
		// If the filter is all or the filter matches the role
		if strings.ToLower(filter) == "all" || strings.ToLower(node.Role) == strings.ToLower(filter) {
			nodeIPs = append(nodeIPs, net.ParseIP(node.NodeIPAddr))
		}
	}

	return nodeIPs
}
