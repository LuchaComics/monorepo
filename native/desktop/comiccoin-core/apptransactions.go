package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) GetTransactions(address string) ([]*domain.BlockTransaction, error) {
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
