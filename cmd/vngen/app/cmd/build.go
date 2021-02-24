package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"nenvoy.com/pkg/topology"
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
			return
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
	// Create the deployment
	err = topology.BuildFromFile(args[0])
	if err != nil {
		return err
	}

	return nil
}
