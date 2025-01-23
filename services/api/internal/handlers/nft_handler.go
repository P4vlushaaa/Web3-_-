package handlers

import (
	"encoding/json"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"net/http"
	"web3-onlyfans/services/api/internal/neo"
	"web3-onlyfans/services/api/internal/utils"
)

type MintNFTRequest struct {
	TokenID    string `json:"token_id"`
	Name       string `json:"name"`
	BlurredRef string `json:"blurred_ref"`
	FullRef    string `json:"full_ref"`
	Price      int    `json:"price"`
}

func MintNFTHandler(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient, logger utils.Logger, cfg *utils.Config) {
	var req MintNFTRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	nftHash, err := util.Uint160DecodeStringLE(cfg.NftContractHash)
	if err != nil {
		http.Error(w, "invalid NFT hash", http.StatusInternalServerError)
		return
	}

	txHash, err := neoCli.Actor.Call(nftHash, "mint", []any{
		req.TokenID,
		req.Name,
		req.BlurredRef,
		req.FullRef,
		req.Price,
	}, nil)
	if err != nil {
		logger.Errorf("mint call error: %v", err)
		http.Error(w, "mint call error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"txHash":  txHash.StringBE(),
		"tokenID": req.TokenID,
	})
}

func NFTPropertiesHandler(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient, logger utils.Logger, cfg *utils.Config) {
	tokenID := r.URL.Query().Get("token_id")
	if tokenID == "" {
		http.Error(w, "missing token_id", http.StatusBadRequest)
		return
	}

	nftHash, err := util.Uint160DecodeStringLE(cfg.NftContractHash)
	if err != nil {
		http.Error(w, "invalid NFT hash", http.StatusInternalServerError)
		return
	}

	invRes, err := neoCli.Actor.Reader().InvokeCall(nftHash, "properties", []any{tokenID}, nil)
	if err != nil {
		logger.Errorf("invoke call error: %v", err)
		http.Error(w, "invoke call error", http.StatusInternalServerError)
		return
	}

	// Разбираем stackitem.Map
	propsMap, err := utils.DecodeStringMap(invRes.Stack[0])
	if err != nil {
		logger.Errorf("decode error: %v", err)
		http.Error(w, "decode error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"props":  propsMap,
	})
}
