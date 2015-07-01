package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Address string            `json:"address"`
	Port    int               `json:"port"`
	Servers map[string]string `json:"servers"`
}

var cfg config

func readConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&cfg)
	return err
}
