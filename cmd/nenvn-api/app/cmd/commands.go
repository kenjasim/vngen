package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"nenvoy.com/cmd/nenvn-api/app/pkg/actions"
	"nenvoy.com/pkg/utils/handle"

	"github.com/gorilla/mux"
	"nenvoy.com/cmd/nenvn-api/app/pkg/details"
)

func build(w http.ResponseWriter, r *http.Request) {
	// Read the http request body
	b, err := ioutil.ReadAll(r.Body)
	handle.Error(err)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error reading template"))
		return
	}

	err = actions.Build(b)
	handle.Error(err)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error building template"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Successfuly built template"))
}

func start(w http.ResponseWriter, r *http.Request) {
	// Get the variables
	vars := mux.Vars(r)

	if vars["resource"] != "host" && vars["resource"] != "deployment" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can only start host or deployment, not %s", vars["resource"])))
		return
	}

	err := actions.Start(vars["name"], vars["resource"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to start %s %s", vars["resource"], vars["name"])))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Successfuly started %s %s", vars["resource"], vars["name"])))
}

func stop(w http.ResponseWriter, r *http.Request) {
	// Get the variables
	vars := mux.Vars(r)
	if vars["resource"] != "host" && vars["resource"] != "deployment" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can only stop host or deployment, not %s", vars["resource"])))
		return
	}

	err := actions.Stop(vars["name"], vars["resource"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to stop %s %s", vars["resource"], vars["name"])))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Successfuly stopped %s %s", vars["resource"], vars["name"])))
}

func restart(w http.ResponseWriter, r *http.Request) {
	// Get the variables
	vars := mux.Vars(r)

	if vars["resource"] != "host" && vars["resource"] != "deployment" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can only restart host or deployment, not %s", vars["resource"])))
		return
	}

	err := actions.Restart(vars["name"], vars["resource"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to restart %s %s", vars["resource"], vars["name"])))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Successfuly restarted %s %s", vars["resource"], vars["name"])))
}

func destroy(w http.ResponseWriter, r *http.Request) {
	// Get the variables
	vars := mux.Vars(r)

	if vars["resource"] != "host" && vars["resource"] != "deployment" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can only destroy host or deployment, not %s", vars["resource"])))
		return
	}

	err := actions.Destroy(vars["name"], vars["resource"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to destroy %s %s", vars["resource"], vars["name"])))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Successfuly destroyed %s %s", vars["resource"], vars["name"])))
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	resp, err := details.GetHosts()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error collecting details"))
		return
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
		return
	}

	// Write the application type headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write(resp)
}

func getHost(w http.ResponseWriter, r *http.Request) {

	// Get the variables
	vars := mux.Vars(r)
	resp, err := details.GetHost(vars["host"])

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error collecting details"))
		return
	}

	// Write the application type headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write(resp)
}

func getHostIP(w http.ResponseWriter, r *http.Request) {

	// Get the variables
	vars := mux.Vars(r)
	resp, err := details.GetHostIP(vars["host"])

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error collecting details"))
		return
	}

	// Write the application type headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write(resp)
}
