package cmd

import (
	"net/http"

	"nenvoy.com/cmd/nenvn-api/app/pkg/details"
)

func getHosts(w http.ResponseWriter, r *http.Request) {
	resp, err := details.GetHosts()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error collecting details"))
	}

	// Write the application type headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write(resp)
}

func getNetworks(w http.ResponseWriter, r *http.Request) {
	resp, err := details.GetNetworks()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error collecting details"))
	}

	// Write the application type headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write(resp)
}
