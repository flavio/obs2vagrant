package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Tests for the findBoxJSONFile function.

func TestFindBoxJSONFileFail(t *testing.T) {
	_, err := findBoxJSONFile("http://", "name")
	if err == nil {
		t.Fatal("It should've failed")
	}
}

func TestFindBoxJSONFileEmpty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader := strings.NewReader("")
		io.Copy(w, reader)
	}))
	defer ts.Close()

	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	_, err := findBoxJSONFile(ts.URL, "name")
	if err == nil || err.Error != "Cannot find box" || err.Code != 404 {
		t.Fatal("Wrong error")
	}
	if !strings.Contains(buffer.String(), "Cannot find box inside of") {
		t.Fatal("Wrong logged message")
	}
}

func TestFindBoxJSONFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader := strings.NewReader(`<href="name.x86_64-1.12.1-Build0.1.json">`)
		io.Copy(w, reader)
	}))
	defer ts.Close()

	result, err := findBoxJSONFile(ts.URL, "name")
	if err != nil {
		t.Fatal("It should've gone correctly")
	}
	if result != "name.x86_64-1.12.1-Build0.1.json" {
		t.Fatal("Wring JSON name")
	}
}

// Tests for the getBoxJSON function.

func TestGetBoxJSONFail(t *testing.T) {
	_, err := getBoxJSON("http://")
	if err == nil {
		t.Fatal("It should've failed")
	}
}

func TestGetBoxJSONMalformed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader := strings.NewReader("invalid json is invalid")
		io.Copy(w, reader)
	}))
	defer ts.Close()

	_, err := getBoxJSON(ts.URL)
	msg := "invalid character 'i' looking for beginning of value"
	if err == nil || err.Error != msg || err.Code != 500 {
		t.Fatal("Wrong error!")
	}
}

func TestGetBoxJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("data/box.json")
		if err != nil {
			fmt.Fprintln(w, "FAIL!")
			return
		}
		io.Copy(w, file)
		file.Close()
	}))
	defer ts.Close()

	box, err := getBoxJSON(ts.URL)
	if err != nil {
		t.Fatal("It should return no errors")
	}

	if box.Description != "Base SLES12 btrfs box for testing docker" {
		t.Fatal("Unexpected value for the description")
	}
	if box.Name != "Base-SLES12-btrfs" {
		t.Fatal("Unexpected value for the name")
	}
	if len(box.Versions) != 1 {
		t.Fatal("Wrong number of versions")
	}
	if box.Versions[0].Version != "1.12.1" {
		t.Fatal("Wrong box version")
	}
	if len(box.Versions[0].Providers) != 1 {
		t.Fatal("Wrong number of providers")
	}
	url := "Base-SLES12-btrfs.x86_64-1.12.1.libvirt-Build2.1.box"
	if box.Versions[0].Providers[0].Url != url {
		t.Fatal("Wrong value for the box provider url")
	}
	if box.Versions[0].Providers[0].Name != "libvirt" {
		t.Fatal("Wrong value for the box provider name")
	}
}
