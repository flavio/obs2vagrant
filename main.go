// This is a simple web application that makes possible to use Open Build
// Service as a simple Vagrant image catalog.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// Returns the "Not found" response.
func notFound(w http.ResponseWriter, req *http.Request) {
	writeError(w, errorResponse{
		Error: "Not Found",
		Code:  http.StatusNotFound,
	})
}

func main() {
	var configFile string
	const defaultConfigFile = "obs2vagrant.json"
	flag.StringVar(&configFile, "c", defaultConfigFile, "configuration file")
	flag.Parse()

	err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Error while parsing configuration file: %s", err)
	}

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/{server}/{project}/{repo}/{box}.json", boxHandler).
		Methods("GET")
	n.UseHandler(r)

	listenOn := fmt.Sprintf("%v:%v", cfg.Address, cfg.Port)
	n.Run(listenOn)
}
