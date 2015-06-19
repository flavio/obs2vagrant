package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func boxHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	serverName := vars["server"]
	project := vars["project"]
	repository := vars["repo"]
	box := vars["box"]

	server := config.Servers[serverName]

	if server == "" {
		log.Printf("ERROR: Cannot find %s server inside of configuration file\n", server)
		writeError(w, errorResponse{
			Error: "Bad Request: Unknown server",
			Code:  http.StatusBadRequest,
		})
		return
	}

	jsonFileName, appErr := findBoxJSONFile(server, project, repository, box)
	if appErr != nil {
		appErr.Write(w)
		return
	}

	jsonBox, appErr := getBoxJSON(server, project, repository, jsonFileName)
	if appErr != nil {
		appErr.Write(w)
		return
	}

	for i := range jsonBox.Versions {
		version := &jsonBox.Versions[i]
		for j := range version.Providers {
			provider := &version.Providers[j]
			provider.Url = server + strings.Replace(project, ":", ":/", -1) + "/" + repository + "/" + provider.Url
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
