package neo

import (
	"context"
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/actor"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

type NeoClient struct {
	Actor   *actor.Actor
	RPC     *rpcclient.Client
	Context context.Context
}

func NewNeoClient(rpcEndpoint, walletPath, walletPass string) (*NeoClient, error) {
	ctx := context.Background()

	rpcCli, err := rpcclient.New(ctx, rpcEndpoint, rpcclient.Options{})
	if err != nil {
		return nil, fmt.Errorf("rpcclient new: %w", err)
	}
	w, err := wallet.NewWalletFromFile(walletPath)
	if err != nil {
		return nil, fmt.Errorf("wallet load: %w", err)
	}
	acc := w.GetAccount(w.GetChangeAddress())
	err = acc.Decrypt(walletPass, w.Scrypt)
	if err != nil {
		return nil, fmt.Errorf("wallet decrypt: %w", err)
	}
	act, err := actor.NewSimple(rpcCli, acc)
	if err != nil {
		return nil, fmt.Errorf("actor new: %w", err)
	}

	return &NeoClient{
		Actor:   act,
		RPC:     rpcCli,
		Context: ctx,
	}, nil
}
