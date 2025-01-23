package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "web3-onlyfans/services/frostfs-uploader/internal/utils"
    "web3-onlyfans/services/frostfs-uploader/internal/handlers"
)

func main() {
    cfgPath := os.Getenv("FROSTFS_UPLOADER_CONFIG")
    if cfgPath == "" {
        cfgPath = "./config.yaml"
    }
    cfg, err := utils.LoadConfig(cfgPath)
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }
    logger := utils.NewLogger(cfg.LogLevel)

    r := mux.NewRouter()
    r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        handlers.UploadHandler(w, r, cfg, logger)
    }).Methods("POST")

    addr := cfg.ListenAddr
    logger.Infof("frostfs-uploader on %s", addr)
    log.Fatal(http.ListenAndServe(addr, r))
}
