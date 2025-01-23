package wallet

import (
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

// Пример – если нужно создавать/импортировать кошельки из seed-фразы.
// Здесь упрощённая заготовка.
func CreateNewWallet(path, password string) error {
	w, err := wallet.NewWallet(path)
	if err != nil {
		return fmt.Errorf("wallet new: %w", err)
	}
	// ...
	return nil
}
