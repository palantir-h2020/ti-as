package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config struct for webapp config
type Config struct {
	// Describing service properties
	Server struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		Host string `yaml:"server_host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port string `yaml:"server_port"`
	} `yaml:"server"`

	// Describing properties of database that will be used
	Database struct {
		// What database will be used. At the moment only redis is supported.
		Instance string `yaml:"db_instance"`

		// In which host is db deployed
		Host string `yaml:"db_host"`

		// In which port is db listening to
		Port int `yaml:"db_port"`

		// Username to connect
		Username string `yaml:"db_username"`

		// Password to connect
		Password string `yaml:"db_password"`

		// Name of the database, that will be used
		Name string `yaml:"db_name"`
	} `yaml:"database"`
}

func (c *Config) LoadConf(filename string) *Config {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
