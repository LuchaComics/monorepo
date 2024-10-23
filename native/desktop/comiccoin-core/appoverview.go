package main

import (
	"fmt"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) GetTotalCoins(address string) (uint64, error) {
	// Defensive code
	if address == "" {
		return 0, fmt.Errorf("failed because: address is null: %v", address)
	}

	//TODO: Impl.

	return 0, nil
}

func (a *App) GetTotalTokens(address string) (uint64, error) {
	// Defensive code
	if address == "" {
		return 0, fmt.Errorf("failed because: address is null: %v", address)
	}

	//TODO: Impl.

	return 0, nil
}

func (a *App) GetRecentTransactions(address string) ([]*domain.Transaction, error) {
	// Defensive code
	if address == "" {
		return make([]*domain.Transaction, 0), fmt.Errorf("failed because: address is null: %v", address)
	}

	//TODO: Impl.

	return make([]*domain.Transaction, 0), nil
}
