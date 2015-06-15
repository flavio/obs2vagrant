package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func findBoxJSONFile(server string, project string, repository string, name string) (string, *appError) {
	pattern := "href=\"(" + name + "[\\w\\d-.]+-Build[\\w\\d-.]+\\.json)\">"
	re := regexp.MustCompile(pattern)

	indexUrl := server + strings.Replace(project+"/"+repository, ":", ":/", -1)
	resp, err := http.Get(indexUrl)
	if err != nil {
		return "", &appError{err, 500}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "",
			&appError{fmt.Errorf("GET %s failed with %s", indexUrl, resp.Status), resp.StatusCode}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", &appError{err, 500}
	}

	matches := re.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		log.Printf("Cannot find box inside of %s", body)
		return "", &appError{fmt.Errorf("Cannot find box"), 404}
	} else if len(matches) > 2 {
		log.Printf("Found more than 2 matches: %+v", matches)
		return "", &appError{fmt.Errorf("Found multiple matches"), 400}
	} else {
		return matches[0][1], nil
	}
}

func getBoxJSON(server string, project string, repository string, jsonFile string) (BoxJSON, *appError) {
	boxJSON := BoxJSON{}

	indexUrl := server + strings.Replace(project+"/"+repository, ":", ":/", -1) + "/" + jsonFile
	resp, err := http.Get(indexUrl)
	if err != nil {
		return boxJSON, &appError{err, 500}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return boxJSON,
			&appError{fmt.Errorf("GET %s failed with %s", indexUrl, resp.Status), resp.StatusCode}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return boxJSON, &appError{err, 500}
	}

	err = json.Unmarshal(body, &boxJSON)
	if err != nil {
		return boxJSON, &appError{err, 500}
	} else {
		return boxJSON, nil
	}
}
