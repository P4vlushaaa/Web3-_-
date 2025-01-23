package main

import (
	"log"
	"web3-onlyfans/services/indexer/internal/subscriber"
	"web3-onlyfans/services/indexer/internal/utils"
)

func main() {
	cfg, err := utils.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger := utils.NewLogger(cfg.LogLevel)

	sub, err := subscriber.NewBlockSubscriber(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to create subscriber: %v", err)
	}

	logger.Infof("Indexer started, polling blocks from %s", cfg.NeoRPC)
	sub.Start()
}
