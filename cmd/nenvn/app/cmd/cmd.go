package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"nenvoy.com/pkg/utils/handle"
	"nenvoy.com/pkg/utils/root"
)

var (
	// This is used for config file
	cfgFile string

	// Setup the initial nenadm command structure
	baseCmd = &cobra.Command{
		Use:   "nenvn",
		Short: "nenvn: CLI to deploy custom virtualised networks using QEMU/KVM",
		Long:  `nenvn: CLI to deploy custom virtualised networks using QEMU/KVM`,
	}
)

// Execute executes the nenadm command.
func Execute() error {
	// Check running as root
	root, err := root.AsRoot()

	if err != nil {
		handle.Error(err)
	}

	if !root {
		handle.Error(errors.New("permission error: root required"))
		os.Exit(1)
	}
	return baseCmd.Execute()
}

func init() {
	// cobra.OnInitialize(readConfig)
	// baseCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "cluster config file (default is config.yaml)")
}

// Reads the config file
// func readConfig() {
// 	// err := config.ReadConfig(&cfgFile)
// 	// if err != nil {
// 	// 	printing.PrintWarning(err.Error())
// 	// }
// }
