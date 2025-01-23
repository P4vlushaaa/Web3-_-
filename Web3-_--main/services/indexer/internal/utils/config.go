package utils

import (
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NeoRPC              string `yaml:"neoRPC"`
	LogLevel            string `yaml:"logLevel"`
	DBDsn               string `yaml:"dbDsn"`
	PollIntervalSeconds int    `yaml:"pollIntervalSeconds"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	// Переопределения через ENV
	if e := os.Getenv("INDEXER_POLL_INTERVAL"); e != "" {
		if val, err := strconv.Atoi(e); err == nil {
			c.PollIntervalSeconds = val
		}
	}

	return &c, nil
}
