package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

func (a *App) GetTransactions(address string) ([]*domain.BlockTransaction, error) {
	addr := common.HexToAddress(strings.ToLower(address))

	// Defensive code
	if address == "" {
		return make([]*domain.BlockTransaction, 0), fmt.Errorf("failed because: address is null: %v", address)
	}

	txs, err := a.listBlockTransactionsByAddressService.Execute(a.ctx, &addr)
	if err != nil {
		a.logger.Error("Failed getting account balance", slog.Any("error", err))
		return nil, err
	}

	return txs, nil
}

func (a *App) GetBlockDataByBlockTransactionTimestamp(blockDataBlockTransactionTimestamp uint64) (*domain.BlockData, error) {
	return a.getByBlockTransactionTimestampService.Execute(a.ctx, blockDataBlockTransactionTimestamp)
}
