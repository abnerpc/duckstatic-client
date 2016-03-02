package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration has the duckstatic-server url and the user credentials to connect
type Configuration struct {
	ServerURL string
	Username  string
	Password  string
}

const defaultConfigFileName string = "ds-client.conf"

// Config has the current configuration
var Config *Configuration

// LoadConfigurationFrom loads the configuration from path and returns a
// Configuration pointer
func LoadConfigurationFrom(path string) error {

	file, err := os.Open(filepath.Join(path, defaultConfigFileName))
	if err != nil {
		return CreateNewConfig(path)
	}

	decoder := json.NewDecoder(file)
	Config = &Configuration{}
	return decoder.Decode(Config)
}

// CreateNewConfig creates config file with user info
func CreateNewConfig(path string) error {
	// TODO
	return nil
}

// WriteConfiguration saves the current Configuration state to configFilePath
func WriteConfiguration(filePath string) error {
	configJSON, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(filePath, configJSON, 0644)
	}
	return err
}
