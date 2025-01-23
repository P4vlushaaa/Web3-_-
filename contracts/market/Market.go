package contract

import (
	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
	"github.com/nspcc-dev/neo-go/pkg/interop/iterator"
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/std"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const (
	ownerKey   = "o"
	tokenKey   = "t"
	nftKey     = "n"
	sellPrefix = "forSale_"
)

type MarketDeployData struct {
	Admin interop.Hash160
	Token interop.Hash160
	Nft   interop.Hash160
}

func _deploy(data any, isUpdate bool) {
	if isUpdate {
		return
	}

	initData := data.(MarketDeployData)
	admin := initData.Admin
	if len(admin) != 20 {
		panic("invalid admin")
	}

	ctx := storage.GetContext()
	storage.Put(ctx, []byte(ownerKey), admin)
	storage.Put(ctx, []byte(tokenKey), initData.Token)
	storage.Put(ctx, []byte(nftKey), initData.Nft)
}


func OnNEP11Payment(from interop.Hash160, amount int, tokenId []byte, data any) {
	if amount != 1 {
		panic("Must be non-fungible (amount=1)")
	}
	ctx := storage.GetContext()
	callingHash := runtime.GetCallingScriptHash()
	nftHash := storage.Get(ctx, []byte(nftKey)).(interop.Hash160)
	if !std.Equals(callingHash, nftHash) {
		panic("NFT not recognized")
	}

	// Выставляем NFT на продажу, если владелец хочет
	// Цена берется из свойств NFT (fields: price)
	saleKey := append([]byte(sellPrefix), tokenId...)
	storage.Put(ctx, saleKey, []byte{1})
}


func OnNEP17Payment(from interop.Hash160, amount int, data any) {
	ctx := storage.GetContext()
	callingHash := runtime.GetCallingScriptHash()
	myToken := storage.Get(ctx, []byte(tokenKey)).(interop.Hash160)
	if !std.Equals(callingHash, myToken) {
		panic("invalid token")
	}

	tokenIdBytes, ok := data.([]byte)
	if !ok {
		panic("no tokenId in data")
	}
	tokenId := string(tokenIdBytes)

	saleKey := append([]byte(sellPrefix), tokenIdBytes...)
	if storage.Get(ctx, saleKey) == nil {
		panic("NFT not on sale")
	}

	
	nftHash := storage.Get(ctx, []byte(nftKey)).(interop.Hash160)
	props := contract.Call(nftHash, "properties", contract.ReadStates, tokenId).(map[string]any)
	// price хранится как string => нужно конвертировать
	price := std.Atoi(props["price"].(string))
	if amount < price {
		panic("insufficient payment")
	}

	
	storage.Delete(ctx, saleKey)
	// передаем NFT покупателю
	contract.Call(nftHash, "transfer", contract.All, from, tokenId, nil)
}


func List() []string {
	ctx := storage.GetContext()
	iter := storage.Find(ctx, []byte(sellPrefix), storage.KeysOnly)
	var result []string
	for iterator.Next(iter) {
		key := iterator.Key(iter).([]byte)
		tokenId := key[len(sellPrefix):]
		result = append(result, string(tokenId))
	}
	return result
}


func TransferTokens(to interop.Hash160, amount int) {
	ctx := storage.GetContext()
	owner := storage.Get(ctx, []byte(ownerKey)).(interop.Hash160)
	if !runtime.CheckWitness(owner) {
		panic("not an owner")
	}
	myToken := storage.Get(ctx, []byte(tokenKey)).(interop.Hash160)
	contract.Call(myToken, "transfer", contract.All, runtime.GetExecutingScriptHash(), to, amount, nil)
}
