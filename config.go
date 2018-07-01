package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Directories []string
	Database    []string
}

func getConfig() *Configuration {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &configuration
}

func getDirectories() []string {
	config := getConfig()
	return config.Directories
}
