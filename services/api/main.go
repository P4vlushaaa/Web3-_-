package main

import (
	"log"
	"net/http"
	"os"

	"myproject/services/api/internal/handlers"
	"myproject/services/api/internal/neo"

	"github.com/gorilla/mux"
)

func main() {
	rpcEndpoint := os.Getenv("NEO_RPC_ENDPOINT") // "http://localhost:20332"
	walletPath := os.Getenv("NEO_WALLET_PATH")
	walletPass := os.Getenv("NEO_WALLET_PASS")

	// Инициализируем клиент NEO
	neoCli, err := neo.NewNeoClient(rpcEndpoint, walletPath, walletPass)
	if err != nil {
		log.Fatalf("failed to init neo client: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/nft/mint", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMintNFT(w, r, neoCli)
	}).Methods("POST")

	r.HandleFunc("/nft/market-list", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleListMarket(w, r, neoCli)
	}).Methods("GET")

	// и т.д.

	log.Println("Starting API on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
