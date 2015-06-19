package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// This is the error to be used when interacting with clients.
type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// Concatenate this response by marshalling it into JSON.
func (r errorResponse) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// Log this error and write it on the given ResponseWriter too in a JSON
// format.
func (r errorResponse) Write(w http.ResponseWriter) {
	log.Printf("ERROR: %s\n", r.Error)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	fmt.Fprint(w, r)
}

// Write the given error to the writer.
func writeError(w http.ResponseWriter, err errorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	fmt.Fprint(w, err)
}
