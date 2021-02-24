package cmd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"nenvoy.com/cmd/vngen/app/pkg/api"
	"nenvoy.com/pkg/utils/handle"
	"nenvoy.com/pkg/utils/printing"
)

func init() {
	// Api port flag
	apiCmd.PersistentFlags().IntVarP(&port, "port", "p", 8000, "Specify the port number to be used for the rest api")

	baseCmd.AddCommand(apiCmd)
}

var (
	port int

	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run the rest api",
		Long:  `Run the rest api`,
		Run: func(cmd *cobra.Command, args []string) {

			printing.PrintInfo(fmt.Sprintf("Running rest api on http://localhost:%d", port))
			// Distribute binaries and handle the setup

			r := mux.NewRouter()

			// Handle the building of the deployment
			r.HandleFunc("/build", api.Build).Methods("PUT")

			// Handle the starting of the deployment or host
			r.HandleFunc("/start/{resource}/{name}", api.Start).Methods("POST")

			// Handle the stopping of the deployment or host
			r.HandleFunc("/stop/{resource}/{name}", api.Stop).Methods("POST")

			// Handle the restarting of the deployment or host
			r.HandleFunc("/restart/{resource}/{name}", api.Restart).Methods("POST")

			// Handle the destroying of the deployment or host
			r.HandleFunc("/destroy/{resource}/{name}", api.Destroy).Methods("POST")

			// Handle the getting of the host details
			r.HandleFunc("/hosts", api.GetHosts)

			// Handle the getting of the network details
			r.HandleFunc("/networks", api.GetNetworks)

			// Get the IP of a particular host
			r.HandleFunc("/details/{host}", api.GetHost)

			// Get the IP of a particular host
			r.HandleFunc("/details/{host}/ipv4", api.GetHostIP)

			handle.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
		},
	}
)
