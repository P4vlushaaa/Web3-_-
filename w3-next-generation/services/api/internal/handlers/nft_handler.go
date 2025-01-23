package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"os"

	"myproject/services/api/internal/neo"

	"github.com/nspcc-dev/neo-go/pkg/util"
)

// Пример: /nft/mint?token_id=myPost1&blurred=OID1&full=OID2&price=100
func HandleMintNFT(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient) {
	// Получаем параметры
	tokenId := r.URL.Query().Get("token_id")
	blurred := r.URL.Query().Get("blurred")
	full := r.URL.Query().Get("full")
	priceStr := r.URL.Query().Get("price")
	if tokenId == "" || blurred == "" || full == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	// Конвертируем price
	// (в примере NEP-17 без decimals, так что не умножаем на 10^дробность)
	price := big.NewInt(0)
	price.SetString(priceStr, 10)

	// Хэш NFT-контракта из ENV
	nftHashStr := os.Getenv("NFT_CONTRACT_HASH")
	nftHash, err := util.Uint160DecodeStringLE(nftHashStr)
	if err != nil {
		http.Error(w, "invalid NFT hash", http.StatusInternalServerError)
		return
	}

	// Вызываем метод "mint"
	txHash, err := neoCli.Actor.Call(nftHash, "mint", []any{
		tokenId,  // str
		"My NFT", // name (можно отдельно передавать)
		blurred,  // blurred_ref
		full,     // full_ref
		price.Int64(),
	}, nil)
	if err != nil {
		http.Error(w, "mint call error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"tx_hash": txHash.StringBE(),
	})
}

func HandleListMarket(w http.ResponseWriter, r *http.Request, neoCli *neo.NeoClient) {
	marketHashStr := os.Getenv("MARKET_CONTRACT_HASH")
	marketHash, err := util.Uint160DecodeStringLE(marketHashStr)
	if err != nil {
		http.Error(w, "invalid MARKET hash", http.StatusInternalServerError)
		return
	}

	// Invoke без отправки транзакции (чтение из стейта)
	res, err := neoCli.Actor.Reader().InvokeCall(marketHash, "List", []any{}, nil)
	if err != nil {
		http.Error(w, "List call error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// res.Value — stackitem
	// Для простоты декодируем
	tokenIDs, err := res.Stack[0].TryArray()
	if err != nil {
		http.Error(w, "decoding array error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	list := []string{}
	for _, t := range tokenIDs {
		bs, err := t.TryBytes()
		if err == nil {
			list = append(list, string(bs))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"onSale": list,
	})
}
