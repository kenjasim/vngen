package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer - Starts the server with all the routing in place
func StartServer() {
	r := mux.NewRouter()

	// Handle the building of the deployment
	r.HandleFunc("/build", build).Methods("PUT")

	// Handle the starting of the deployment or host
	r.HandleFunc("/start/{resource}/{name}", start).Methods("POST")

	// Handle the stopping of the deployment or host
	r.HandleFunc("/stop/{resource}/{name}", stop).Methods("POST")

	// Handle the restarting of the deployment or host
	r.HandleFunc("/restart/{resource}/{name}", restart).Methods("POST")

	// Handle the destroying of the deployment or host
	r.HandleFunc("/destroy/{resource}/{name}", destroy).Methods("POST")

	// Handle the getting of the host details
	r.HandleFunc("/hosts", getHosts)

	// Handle the getting of the network details
	r.HandleFunc("/networks", getNetworks)

	// Get the IP of a particular host
	r.HandleFunc("/details/{host}", getHost)

	// Get the IP of a particular host
	r.HandleFunc("/details/{host}/ipv4", getHostIP)

	http.ListenAndServe(":8000", r)

}
