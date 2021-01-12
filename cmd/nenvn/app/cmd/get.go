package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"

	"github.com/spf13/cobra"
	"nenvoy.com/pkg/utils/handle"
	"nenvoy.com/pkg/utils/printing"
)

func init() {
	baseCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get <hosts|networks|deployments>",
	Short: "Return objects from the database",
	Long:  `Return objects from the database`,
	Run: func(cmd *cobra.Command, args []string) {

		//Sepcify the username for connecting
		if len(args) != 1 {
			handle.Error(errors.New("Need to specify template file, see help for more details"))
		}

		printing.PrintInfo(fmt.Sprintf("Getting %s", args[0]))
		// Get the hosts
		if args[0] == "hosts" {
			handle.Error(getHosts())
		} else if args[0] == "networks" {
			handle.Error(getNetworks())
		}

	},
}

func getHosts() (err error) {

	// Get all the hosts from the database
	hosts, err := host.GetHosts()
	if err != nil {
		return err
	}

	// Create the table and print the hosts
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tRAM\tCPU\tStorage\tDeployment\t")

	for _, host := range hosts {
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%d\n", host.Name, host.RAM, host.CPUs, host.HDSpace, host.DeploymentID)
	}
	w.Flush()

	return nil
}

func getNetworks() (err error) {
	// Get all the hosts from the database
	networks, err := network.GetNetworks()
	if err != nil {
		return err
	}

	// Create the table and print the hosts
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tType\tIP\tDHCP Range\tDeploymentID\t")

	for _, network := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", network.Name, network.Type, network.IP, fmt.Sprintf("%s-%s", network.DHCPLower, network.DHCPUpper), network.DeploymentID)
	}
	w.Flush()

	return nil
}
