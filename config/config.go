package config

import (
	"encoding/json"
	"fmt"
	"goserver/utils/gsclient"
	"goserver/utils/gslog"
	"os"
)

var Conf Config

type ClientKey string

const (
	MOCK_CLIENT ClientKey = "MockClient"
)

var clientKeys = []ClientKey{
	MOCK_CLIENT,
}

type Config struct {
	Port         int 	                      `json:"port"`
	Basepath     string                       `json:"basepath"`
	Logger       *gslog.LoggerConfig          `json:"logger"`
	LogFile      gslog.LogFileConfig          `json:"log_file"`
	TokenClients []gsclient.ClientConfig      `json:"token_clients"`
	TokenSources []gsclient.TokenSourceConfig `json:"token_sources"`
	Clients      []gsclient.ClientConfig      `json:"clients"`
}

func (c *Config) validateClients() error {
	for _, requiredKey := range clientKeys {
		found := false
		for _, client := range c.Clients {
			if client.Key == string(requiredKey) {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("couldn't find required client with key %s in configuration file", string(requiredKey))
		}
	}
	return nil
}

func LoadConfiguration(f string) error {

	// initialize default values
	var config Config

	// open given file
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	// parse file content as JSON
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&config)
	if err != nil {
		return err
	}

	if config.Logger != nil {
		gslog.ConfigureLog(*config.Logger)
	}
	gslog.ConfigureLogFile(config.LogFile)
	
	err = gsclient.ConfigureClients(config.TokenClients)
	if err != nil {
		return err
	}

	err = gsclient.ConfigureTokenSources(config.TokenSources)
	if err != nil {
		return err
	}

	err = gsclient.ConfigureClients(config.Clients)
	if err != nil {
		return err
	}

	return config.validateClients()
}