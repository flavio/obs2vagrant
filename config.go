package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address string            `json:"address"`
	Port    int               `json:"port"`
	Servers map[string]string `json:"servers"`
}

var config Config

func readConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&config)
	return err
}
