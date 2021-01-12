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
	baseCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy <host|deployment> [name]",
	Short: "Undefines hosts and networks in a deployment or a single host",
	Long:  `Undefines hosts and networks in a deployment or a single host`,
	Run: func(cmd *cobra.Command, args []string) {

		//Sepcify the username for connecting
		if len(args) != 2 {
			handle.Error(errors.New("Need to specify deployment or host, see help for more details"))
		}

		printing.PrintInfo(fmt.Sprintf("Destroying %s %s", args[0], args[1]))
		// Get the hosts
		if args[0] == "host" {
			handle.Error(topology.DestroyHost(args[1]))
		} else if args[0] == "deployment" {
			handle.Error(topology.DestroyDeployment(args[1]))
		}

	},
}
