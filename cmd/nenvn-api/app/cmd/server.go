package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer - Starts the server with all the routing in place
func StartServer() {
	r := mux.NewRouter()

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
