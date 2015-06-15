package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Address string               `json:"address"`
	Port    int                  `json:"port"`
	Servers map[string]OBSServer `json:"servers"`
}

type OBSServer struct {
	ApiUrl      string `json:"api_url"`
	DownloadUrl string `json:"download_url"`
	User        string `json:"user"`
	Password    string `json:"password"`
}

func readConfig(config *Config, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
