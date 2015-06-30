package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	url := buildUrl("server/", "project", "repo", "")
	if url != "server/project/repo" {
		t.Fatal("Wrong value for url")
	}

	url = buildUrl("server/", "project", "repo", "extra")
	if url != "server/project/repo/extra" {
		t.Fatal("Wrong value for url")
	}

	url = buildUrl("server/", "project:with:colons", "repo", "extra")
	if url != "server/project:/with:/colons/repo/extra" {
		t.Fatal("Wrong value for url")
	}
}

func TestGetRequestFail(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)

	_, err := getRequest("http://")
	if err == nil || err.Error != "Internal Server Error" || err.Code != 500 {
		t.Fatal("Wrong error")
	}

	msg := "Get http://: http: no Host in request URL\n"
	if !strings.Contains(buffer.String(), msg) {
		t.Fatal("GET request has failed for an unexpected reason")
	}
}

func TestGetRequestNotOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	_, err := getRequest(ts.URL)
	msg := fmt.Sprintf("GET %v failed with 404 Not Found", ts.URL)
	if err.Error != msg || err.Code != 404 {
		t.Fatal("GET request has failed for an unexpected reason")
	}
}

func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}))
	defer ts.Close()

	body, err := getRequest(ts.URL)
	if err != nil {
		t.Fatal("Something went wrong!")
	}
	if string(body) != "OK" {
		t.Fatal("Wrong message!")
	}
}
