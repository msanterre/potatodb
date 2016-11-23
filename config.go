package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigPath = "./config.yaml"
)

type Config struct {
	HttpAddr   string `yaml:"http_addr"`
	SocketAddr string `yaml:"socket_addr"`
	DBFilepath string `yaml:"db_filepath"`
	SaveFreq   int    `yaml:"save_freq"`
}

func LoadConfig() (*Config, error) {

	// Set defaults here
	configs := &Config{
		HttpAddr:   "localhost:5050",
		SocketAddr: "localhost:5051",
		DBFilepath: "data.json",
		SaveFreq:   5,
	}

	content, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &configs)

	return configs, err
}
