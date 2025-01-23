package utils

import (
	"os"

	"io/ioutil"
)

type Config struct {
	ListenAddr         string `yaml:"listenAddr"`
	LogLevel           string `yaml:"logLevel"`
	NeoRPC             string `yaml:"neoRPC"`
	WalletPath         string `yaml:"walletPath"`
	WalletPass         string `yaml:"walletPass"`
	NftContractHash    string `yaml:"nftContractHash"`
	MarketContractHash string `yaml:"marketContractHash"`
	TokenContractHash  string `yaml:"tokenContractHash"`
}

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	if envListen := os.Getenv("API_LISTEN_ADDR"); envListen != "" {
		cfg.ListenAddr = envListen
	}

	return &cfg, nil
}
