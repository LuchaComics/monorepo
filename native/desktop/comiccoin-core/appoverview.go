package main

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) GetTotalCoins(address string) (uint64, error) {
	addr := common.HexToAddress(address)

	// Defensive code
	if address == "" {
		return 0, fmt.Errorf("failed because: address is null: %v", address)
	}

	account, err := a.getAccountService.Execute(&addr)
	if err != nil {
		a.logger.Error("Failed getting account balance", slog.Any("error", err))
		return 0, err
	}

	// Defensive code
	if account == nil {
		err := fmt.Errorf("Failed getting account because D.N.E. at address: %v", addr)
		a.logger.Error("Failed getting account balance", slog.Any("error", err))
		return 0, err
	}

	return account.Balance, nil
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
