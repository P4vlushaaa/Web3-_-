package wallet

import (
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

func CreateNewWallet(path, password string) error {
	w, err := wallet.NewWallet(path)
	if err != nil {
		return fmt.Errorf("wallet new: %w", err)
	}
	return nil
}
