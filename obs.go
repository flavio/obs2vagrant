package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//<directory>
//  <entry name="Base-SLES12-btrfs.x86_64-1.12.1.libvirt-Build2.1.box" />
//  <entry name="Base-SLES12-btrfs.x86_64-1.12.1.libvirt-Build2.1.json" />
//  <entry name="Base-SLES12-ext4.x86_64-1.12.1.libvirt-Build11.1.box" />
//  <entry name="Base-SLES12-ext4.x86_64-1.12.1.libvirt-Build11.1.json" />
//</directory>

type Entry struct {
	XMLName xml.Name `xml:entry`
	Name    string   `xml:"name,attr"`
}

type OBSBinaries struct {
	XMLName xml.Name `xml:directory`
	Entries []Entry  `xml:"entry"`
}

//{
//   "description" : "            Base SLES12 btrfs box for testing docker        ",
//   "versions" : [
//      {
//         "version" : "1.12.1",
//         "providers" : [
//            {
//               "url" : "Base-SLES12-btrfs.x86_64-1.12.1.libvirt-Build2.1.box",
//               "name" : "libvirt"
//            }
//         ]
//      }
//   ],
//   "name" : "Base-SLES12-btrfs"
//}
type BoxJSON struct {
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
	Name        string    `json:"name"`
}

type Version struct {
	Version   string     `json:"version"`
	Providers []Provider `json:"providers"`
}

type Provider struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func obsGetRequest(server OBSServer, request string) ([]byte, *appError) {
	var data []byte
	client := &http.Client{}
	remote_url := server.ApiUrl + request
	req, err := http.NewRequest("GET", remote_url, nil)
	if err != nil {
		return data, &appError{err, 500}
	}
	auth := url.UserPassword(server.User, server.Password)
	req.URL.User = auth

	res, err := client.Do(req)
	if err != nil {
		return data, &appError{err, 500}
	}
	if res.StatusCode != 200 {
		return data,
			&appError{fmt.Errorf("GET %s returned %s", remote_url, res.Status), res.StatusCode}
	}

	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return data, &appError{err, 500}
	}
	return data, nil
}

func getPublishedBinaries(server OBSServer, project string, repository string) (OBSBinaries, *appError) {
	binaries := OBSBinaries{}
	data, appErr := obsGetRequest(server, "/published/"+project+"/"+repository)
	if appErr != nil {
		return binaries, appErr
	}

	err := xml.Unmarshal(data, &binaries)
	if err != nil {
		return binaries, &appError{err, 500}
	} else {
		return binaries, nil
	}
}

func findBox(box string, binaries OBSBinaries) (string, string) {
	var boxFile, jsonFile string
	for _, entry := range binaries.Entries {
		if strings.Contains(entry.Name, box) {
			if strings.HasSuffix(entry.Name, ".box") {
				boxFile = entry.Name
			} else if strings.HasSuffix(entry.Name, ".json") {
				jsonFile = entry.Name
			}
		}
	}

	if boxFile == "" || jsonFile == "" {
		log.Printf("[findBox] looking for box %s inside of %+v\n", box, binaries)
	}

	return boxFile, jsonFile
}

func getBoxJSON(server OBSServer, project string, repository string, jsonFile string) (BoxJSON, *appError) {
	boxJSON := BoxJSON{}
	data, appErr := obsGetRequest(server, "/published/"+project+"/"+repository+"/"+jsonFile)
	if appErr != nil {
		return boxJSON, appErr
	}

	err := json.Unmarshal(data, &boxJSON)
	if err != nil {
		return boxJSON, &appError{err, 500}
	} else {
		return boxJSON, nil
	}
}
