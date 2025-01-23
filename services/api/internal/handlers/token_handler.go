package handlers

import (
	"encoding/json"
	"net/http"
	"web3-onlyfans/services/api/internal/neo"

	"github.com/nspcc-dev/neo-go/pkg/util"
	"web3-onlyfans/services/api/internal/utils"
)

func TokenBalanceHandler(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient, logger utils.Logger, cfg *utils.Config) {
	accountParam := r.URL.Query().Get("account")
	if accountParam == "" {
		accountParam = neoCli.Actor.Sender().StringLE()
	}

	tokenHash, err := util.Uint160DecodeStringLE(cfg.TokenContractHash)
	if err != nil {
		http.Error(w, "invalid token hash", http.StatusInternalServerError)
		return
	}

	accHash, err := util.Uint160DecodeStringLE(accountParam)
	if err != nil {
		http.Error(w, "invalid account param", http.StatusBadRequest)
		return
	}

	invRes, err := neoCli.Actor.Reader().InvokeCall(tokenHash, "balanceOf", []any{accHash}, nil)
	if err != nil {
		logger.Errorf("balanceOf error: %v", err)
		http.Error(w, "balanceOf error", http.StatusInternalServerError)
		return
	}

	balanceI, err := utils.DecodeBigInteger(invRes.Stack[0])
	if err != nil {
		http.Error(w, "decode balance error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"account": accountParam,
		"balance": balanceI.String(),
	})
}
