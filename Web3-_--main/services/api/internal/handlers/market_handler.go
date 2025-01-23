package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"web3-onlyfans/services/api/internal/neo"

	"github.com/nspcc-dev/neo-go/pkg/util"
	"web3-onlyfans/services/api/internal/utils"
)

func MarketListHandler(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient, logger utils.Logger, cfg *utils.Config) {
	marketHash, err := util.Uint160DecodeStringLE(cfg.MarketContractHash)
	if err != nil {
		http.Error(w, "invalid market hash", http.StatusInternalServerError)
		return
	}

	invRes, err := neoCli.Actor.Reader().InvokeCall(marketHash, "List", nil, nil)
	if err != nil {
		logger.Errorf("invoke List error: %v", err)
		http.Error(w, "market list error", http.StatusInternalServerError)
		return
	}

	tokenList, err := utils.DecodeStringArray(invRes.Stack[0])
	if err != nil {
		logger.Errorf("decode array error: %v", err)
		http.Error(w, "decode array error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"tokens": tokenList,
	})
}

type MarketBuyRequest struct {
	TokenID string `json:"token_id"`
	Amount  int64  `json:"amount"`
}

func MarketBuyHandler(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient, logger utils.Logger, cfg *utils.Config) {
	var req MarketBuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	tokenHash, err := util.Uint160DecodeStringLE(cfg.TokenContractHash)
	if err != nil {
		http.Error(w, "invalid token contract hash", http.StatusInternalServerError)
		return
	}

	// Вызываем transfer NEP-17: (from=Actor, to=market, amount=req.Amount, data=req.TokenID)
	dataBytes := []byte(req.TokenID)
	amt := big.NewInt(req.Amount)
	txHash, err := neoCli.Actor.Call(tokenHash, "transfer", []any{
		neoCli.Actor.Sender(),  // from
		cfg.MarketContractHash, // to
		amt,                    // amount
		dataBytes,              // data
	}, nil)
	if err != nil {
		logger.Errorf("market buy transfer error: %v", err)
		http.Error(w, "market buy error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"txHash":  txHash.StringBE(),
		"tokenID": req.TokenID,
		"amount":  req.Amount,
	})
}
