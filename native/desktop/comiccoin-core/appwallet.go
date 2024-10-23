package main

import (
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) ListWallets() ([]*domain.Wallet, error) {
	return a.walletListService.Execute()
}
