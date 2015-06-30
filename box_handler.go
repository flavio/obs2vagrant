package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func boxHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	serverName := vars["server"]
	project := vars["project"]
	repository := vars["repo"]
	box := vars["box"]

	server := cfg.Servers[serverName]

	if server == "" {
		log.Printf("ERROR: Cannot find %s server inside of configuration file\n", server)
		writeError(w, errorResponse{
			Error: "Bad Request: Unknown server",
			Code:  http.StatusBadRequest,
		})
		return
	}

	url := buildUrl(server, project, repository, "")
	jsonFileName, appErr := findBoxJSONFile(url, box)
	if appErr != nil {
		appErr.Write(w)
		return
	}

	url = buildUrl(server, project, repository, jsonFileName)
	jsonBox, appErr := getBoxJSON(url)
	if appErr != nil {
		appErr.Write(w)
		return
	}

	for i := range jsonBox.Versions {
		version := &jsonBox.Versions[i]
		for j := range version.Providers {
			provider := &version.Providers[j]
			provider.Url = buildUrl(server, project, repository, provider.Url)
		}
	}

	jsonData, err := json.Marshal(jsonBox)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		writeError(w, errorResponse{
			Error: "Internal Server Error",
			Code:  http.StatusInternalServerError,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", jsonData)
	}
}
