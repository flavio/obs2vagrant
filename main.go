package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	var configFile string
	const defaultConfigFile = "obs2vagrant.json"
	flag.StringVar(&configFile, "c", defaultConfigFile, "configuration file")
	flag.Parse()

	err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Error while parsing configuration file: %s", err)
	}

	commonHandlers := alice.New(loggingHandler, recoverHandler)

	r := mux.NewRouter()
	r.Handle("/{server}/{project}/{repo}/{box}.json", commonHandlers.ThenFunc(boxHandler))
	http.Handle("/", r)

	listenOn := config.Address + ":" + strconv.FormatInt(int64(config.Port), 10)
	log.Printf("Listening on %s", listenOn)
	http.ListenAndServe(listenOn, nil)
}
