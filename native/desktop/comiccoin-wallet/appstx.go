package main

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) ListAllPendingSignedTransactions() ([]*domain.PendingSignedTransaction, error) {
	pstxs, err := a.pendingSignedTransactionListService.Execute(a.ctx)
	if err != nil {
		a.logger.Error("Failed listing pending stxs by address", slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Fetched pending signed transactions",
		slog.Any("pstxs", pstxs))

	return pstxs, nil
}
