package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func boxHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	serverName := vars["server"]
	project := vars["project"]
	repository := vars["repo"]
	box := vars["box"]

	server := config.Servers[serverName]
	emptyServer := OBSServer{}

	if server == emptyServer {
		log.Printf("ERROR: Cannot find %s server inside of configuration file\n", server)
		http.Error(w, "Error 400: bad request - unknown server", 400)
		return
	}

	binaries, appErr := getPublishedBinaries(server, project, repository)
	if appErr != nil {
		log.Printf("ERROR: %s\n", appErr.Error)
		http.Error(w, http.StatusText(appErr.Code), appErr.Code)
		return
	}
	boxFile, jsonFile := findBox(box, binaries)
	if boxFile == "" || jsonFile == "" {
		log.Println("ERROR: box not found")
		http.NotFound(w, request)
		return
	}

	jsonBox, appErr := getBoxJSON(server, project, repository, jsonFile)
	if appErr != nil {
		log.Printf("ERROR: %s\n", appErr.Error)
		http.Error(w, http.StatusText(appErr.Code), appErr.Code)
		return
	}

	for i := range jsonBox.Versions {
		version := &jsonBox.Versions[i]
		for j := range version.Providers {
			provider := &version.Providers[j]
			provider.Url = server.DownloadUrl + strings.Replace(project, ":", ":/", -1) + "/" + repository + "/" + provider.Url
		}
	}

	jsonData, err := json.Marshal(jsonBox)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", jsonData)
	}
}
