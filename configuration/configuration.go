package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration is a struct representing a single configuration.
type Configuration struct {
	Port string
}

// Config represents the actual configuration.
var Config = new(Configuration)

// Load parses the yml file passed as argument and fills the Config.
func Load(file string) error {
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(conf, &Config)
	if err != nil {
		return err
	}
	return nil
}
