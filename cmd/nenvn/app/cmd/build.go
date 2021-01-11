package cmd

import (
	"errors"

	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"

	"github.com/spf13/cobra"
	"nenvoy.com/pkg/constructor"
	"nenvoy.com/pkg/utils/handle"
	"nenvoy.com/pkg/utils/printing"

	"nenvoy.com/pkg/utils/files"
)

func init() {
	baseCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build <path/to/template>",
	Short: "Build a network from a YAML template file",
	Long:  `Build a network from a YAML template file`,
	Run: func(cmd *cobra.Command, args []string) {

		//Sepcify the username for connecting
		if len(args) != 1 {
			handle.Error(errors.New("Need to specify template file, see help for more details"))
		}

		printing.PrintInfo("Building network from template " + args[0])
		// Distribute binaries and handle the setup
		handle.Error(buildNetwork(args))
	},
}

func buildNetwork(args []string) (err error) {

	// Create relevent directories as needed
	appDir := "/var/lib/nenvn"
	dirs := []string{appDir, appDir + "/machines", appDir + "/images"}
	err = files.CreateDirectories(dirs)
	if err != nil {
		return err
	}
	printing.PrintSuccess("Created application directories")

	// Read in the template file
	netDef, err := constructor.ConvertYAML(args[0])

	// Create the network definitions and create them
	for _, netwr := range netDef.Networks {
		networkDef, err := network.CreateNetworkXML(netwr)
		if err != nil {
			return err
		}

		// Create the network
		err = network.CreateNetwork(networkDef)
		if err != nil {
			return err
		}
	}

	// Create the host definition files and create the hosts
	for _, hst := range netDef.Host {
		hostDef, err := host.CreateHostXML(hst)
		if err != nil {
			return err
		}

		err = host.CreateHost(hostDef)
		if err != nil {
			return err
		}
	}

	return nil
}
