package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) GetTotalCoins(address string) (uint64, error) {
	addr := common.HexToAddress(strings.ToLower(address))

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

func (a *App) GetTotalTokens(address string) (int, error) {
	addr := common.HexToAddress(strings.ToLower(address))

	// Defensive code
	if address == "" {
		return 0, fmt.Errorf("failed because: address is null: %v", address)
	}

	tokCount, err := a.countByOwnerTokenService.Execute(&addr)
	if err != nil {
		a.logger.Error("Failed getting account balance", slog.Any("error", err))
		return 0, err
	}

	return tokCount, nil
}

func (a *App) GetRecentTransactions(address string) ([]*domain.BlockTransaction, error) {
	addr := common.HexToAddress(strings.ToLower(address))

	// Defensive code
	if address == "" {
		return make([]*domain.BlockTransaction, 0), fmt.Errorf("failed because: address is null: %v", address)
	}

	txs, err := a.listRecentBlockTransactionService.Execute(&addr, 5)
	if err != nil {
		a.logger.Error("Failed getting account balance", slog.Any("error", err))
		return nil, err
	}

	return txs, nil
}
