package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "web3-onlyfans/services/api/internal/handlers"
    "web3-onlyfans/services/api/internal/utils"
    "web3-onlyfans/services/api/internal/neo"
)

func main() {
    // Загружаем конфиг (примерно)
    configPath := os.Getenv("API_CONFIG_PATH")
    if configPath == "" {
        configPath = "./config.yaml"
    }
    cfg, err := utils.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Инициализируем логгер
    logger := utils.NewLogger(cfg.LogLevel)

    // Инициализируем NeoClient (RPC + кошелёк)
    neoCli, err := neo.NewNeoClient(cfg.NeoRPC, cfg.WalletPath, cfg.WalletPass)
    if err != nil {
        logger.Fatalf("failed to init neo client: %v", err)
    }

    // Роуты
    r := mux.NewRouter()
    // NFT endpoints
    r.HandleFunc("/nft/mint", func(w http.ResponseWriter, r *http.Request) {
        handlers.MintNFTHandler(w, r, neoCli, logger, cfg)
    }).Methods("POST")

    r.HandleFunc("/nft/properties", func(w http.ResponseWriter, r *http.Request) {
        handlers.NFTPropertiesHandler(w, r, neoCli, logger, cfg)
    }).Methods("GET")

    // Market endpoints
    r.HandleFunc("/market/list", func(w http.ResponseWriter, r *http.Request) {
        handlers.MarketListHandler(w, r, neoCli, logger, cfg)
    }).Methods("GET")

    r.HandleFunc("/market/buy", func(w http.ResponseWriter, r *http.Request) {
        handlers.MarketBuyHandler(w, r, neoCli, logger, cfg)
    }).Methods("POST")

    // Token endpoints (например, посмотреть баланс)
    r.HandleFunc("/token/balance", func(w http.ResponseWriter, r *http.Request) {
        handlers.TokenBalanceHandler(w, r, neoCli, logger, cfg)
    }).Methods("GET")

    // Запуск сервера
    addr := cfg.ListenAddr
    logger.Infof("API starting on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        logger.Fatalf("server error: %v", err)
    }
}
