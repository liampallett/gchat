package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Nick   string `json:"nick"`
	User   string `json:"user"`
	Server string `json:"server"`
	Port   int    `json:"port"`
}

func loadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
