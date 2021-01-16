package cmd

import (
	"errors"
	"fmt"

	"nenvoy.com/pkg/topology"

	"github.com/spf13/cobra"
	"nenvoy.com/pkg/utils/handle"
	"nenvoy.com/pkg/utils/printing"
)

func init() {
	baseCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop <host|deployment> [name]",
	Short: "Stops hosts in a deployment or a single host",
	Long:  `Stops hosts in a deployment or a single host`,
	Run: func(cmd *cobra.Command, args []string) {

		//Sepcify the username for connecting
		if len(args) != 2 {
			handle.Error(errors.New("Need to specify deployment or host, see help for more details"))
			return
		}

		printing.PrintInfo(fmt.Sprintf("Stopping %s %s", args[0], args[1]))
		// Get the hosts
		if args[0] == "host" {
			handle.Error(topology.StopHost(args[1]))
		} else if args[0] == "deployment" {
			handle.Error(topology.StopDeployment(args[1]))
		}

	},
}
