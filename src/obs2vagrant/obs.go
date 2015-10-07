package main

import (
	"encoding/json"
	"log"
	"regexp"
)

type provider struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type version struct {
	Version   string     `json:"version"`
	Providers []provider `json:"providers"`
}

type boxJSON struct {
	Description string    `json:"description"`
	Versions    []version `json:"versions"`
	Name        string    `json:"name"`
}

func findBoxJSONFile(url, name string) (string, *errorResponse) {
	body, err := getRequest(url)
	if err != nil {
		return "", err
	}

	pattern := "href=\"(" + name + "[\\w\\d-.]+-Build[\\w\\d-.]+\\.json)\">"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(body), 1)

	if len(matches) == 0 {
		log.Printf("Cannot find box inside of %s", body)
		return "", &errorResponse{"Cannot find box", 404}
	}
	return matches[0][1], nil
}

func getBoxJSON(url string) (boxJSON, *errorResponse) {
	box := boxJSON{}

	body, err := getRequest(url)
	if err != nil {
		return box, err
	}
	if e := json.Unmarshal(body, &box); e != nil {
		return box, &errorResponse{e.Error(), 500}
	}

	buildRev := getBuildRev(url)
	if buildRev == "" {
		return box, nil
	}

	for i := range box.Versions {
		version := &box.Versions[i]
		// version and build rev must be joined with a '.' otherwise
		// vagrant won't work
		version.Version = version.Version + "." + buildRev
	}
	return box, nil
}

func getBuildRev(url string) string {
	pattern := "[\\w\\d-.]+-Build([\\w\\d-.]+)\\.json"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(url, 1)

	if len(matches) == 0 {
		log.Printf("Cannot find box revision inside of %s", url)
		return ""
	}
	return matches[0][1]
}
