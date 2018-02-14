package database

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var configPath = "postgresConfig.yaml"

//PostgresConfigStr structs yaml configuration
type PostgresConfig struct {
	Dbhost                   string `yaml:"dbhost"`
	Dbport                   string `yaml:"dbport"`
	DbUser                   string `yaml:"dbUser"`
	DbPassword               string `yaml:"dbPassword"`
	DbName                   string `yaml:"dbName"`
	MaxOpenedConnectionsToDb int    `yaml:"maxOpenedConnectionsToDb"`
	MaxIdleConnectionsToDb   int    `yaml:"maxIdleConnectionsToDb"`
	MbConnMaxLifetimeMinutes int    `yaml:"mbConnMaxLifetimeMinutes"`
}

//ReadConfig reads config from file
func ReadConfig() (*PostgresConfig, error) {
	config := &PostgresConfig{}
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
