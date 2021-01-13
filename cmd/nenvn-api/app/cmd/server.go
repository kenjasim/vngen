package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()
	// Handle the getting of the host details
	r.HandleFunc("/details/hosts", getHosts)

	// Handle the getting of the network details
	r.HandleFunc("/details/networks", getNetworks)
	// http.Handle("/", r)
	http.ListenAndServe(":8000", r)

}
