package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"nenvoy.com/pkg/deployment"

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
	Use:   "get <hosts|networks|ips>",
	Short: "Return objects from the database",
	Long:  `Return objects from the database`,
	Run: func(cmd *cobra.Command, args []string) {

		//Sepcify the username for connecting
		if len(args) != 1 {
			handle.Error(errors.New("Need to specify hosts, networks or ips, see help for more details"))
		}

		printing.PrintInfo(fmt.Sprintf("Getting %s", args[0]))
		// Get the hosts
		if args[0] == "hosts" {
			handle.Error(getHosts())
		} else if args[0] == "networks" {
			handle.Error(getNetworks())
		} else if args[0] == "ips" {
			handle.Error(getIPs())
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
	fmt.Fprintln(w, "Name\tVMState\tImage\tRAM\tCPU\tStorage\tDeployment\t")

	for _, host := range hosts {
		// Get the deployment name
		dep, err := deployment.GetDeploymentByID(host.DeploymentID)
		if err != nil {
			return err
		}
		state, err := host.GetHostState()
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\t%s\n", host.Name, state, host.Image, host.RAM, host.CPUs, host.HDSpace, dep.Name)
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
	fmt.Fprintln(w, "Name\tType\tIP\tDHCP Range\tDeployment\t")

	for _, network := range networks {
		// Get the deployment name
		dep, err := deployment.GetDeploymentByID(network.DeploymentID)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", network.Name, network.Type, network.IP, fmt.Sprintf("%s - %s", network.DHCPLower, network.DHCPUpper), dep.Name)
	}
	w.Flush()

	return nil
}

// getIPs - Get the IPs of machines in a deployment
func getIPs() (err error) {
	// Get all the hosts from the database
	hosts, err := host.GetHosts()
	if err != nil {
		return err
	}

	// Create the table and print the hosts
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tInterface\tMacAddr\tIPs\tDeployment\t")

	for _, hst := range hosts {
		dep, err := deployment.GetDeploymentByID(hst.DeploymentID)
		if err != nil {
			return err
		}
		ifaces, err := hst.GetHostIfaces()
		if len(ifaces) == 0 {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", hst.Name, "", "", "", dep.Name)
		}
		for name, iface := range ifaces {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", hst.Name, name, iface[0], strings.Join(iface[1:], ","), dep.Name)
		}
	}
	w.Flush()
	return nil
}
