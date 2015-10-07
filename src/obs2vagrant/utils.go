package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Build a URL as expected by OBS.
func buildUrl(server, project, repo, extra string) string {
	url := server + strings.Replace(project+"/"+repo, ":", ":/", -1)
	if extra != "" {
		url += "/" + extra
	}
	return url
}

// Perform a GET request and return the read body.
func getRequest(url string) ([]byte, *errorResponse) {
	// Check for any errors.
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return []byte{}, &errorResponse{"Internal Server Error", 500}
	}
	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		str := fmt.Sprintf("GET %s failed with %s", url, resp.Status)
		return []byte{}, &errorResponse{str, resp.StatusCode}
	}

	// And read the response itself. Note that `resp.Body` will always be
	// valid, even on error, so we are safe here.
	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return body, nil
}
