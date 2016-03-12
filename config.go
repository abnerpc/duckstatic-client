package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Configuration has the duckstatic-server url and the user credentials to connect
type Configuration struct {
	ServerURL string
	AccessKey string
}

var CurrentUser *user.User

var defaultConfigFilePath string = filepath.Join(CurrentUser.HomeDir, ".config", "ds-client", "config")

// Config has the current configuration
var Config *Configuration

// LoadConfiguration loads the configuration from path and returns a
// Configuration pointer
func LoadConfiguration() error {

	file, err := os.Open(defaultConfigFilePath)
	if err != nil {
		Config = &Configuration{}
		WriteConfiguration()
		return nil
	}

	decoder := json.NewDecoder(file)
	Config = &Configuration{}
	return decoder.Decode(Config)
}

// WriteConfiguration saves the current Configuration state to configFilePath
func WriteConfiguration() error {
	configJSON, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(defaultConfigFilePath, configJSON, 0644)
	}
	return err
}

func (c *Configuration) UpdateServerURL(serverURL string) {
	c.ServerURL = serverURL
	WriteConfiguration()
}
