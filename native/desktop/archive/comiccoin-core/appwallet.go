package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) DefaultWalletAddress() string {
	preferences := PreferencesInstance()
	return preferences.DefaultWalletAddress
}

func (a *App) ListWallets() ([]*domain.Wallet, error) {
	return a.walletListService.Execute()
}

func (a *App) CreateWallet(walletPassword, walletPasswordRepeated, walletLabel string) (string, error) {
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory

	account, err := a.createAccountService.Execute(dataDir, walletPassword, walletPasswordRepeated, walletLabel)
	if err != nil {
		a.logger.Error("failed creating wallet", slog.Any("error", err))
		return "", err
	}
	if account == nil {
		a.logger.Error("failed creating wallet as returned account d.n.e.")
		return "", fmt.Errorf("failed creating wallet: %v", "returned account d.n.e.")
	}

	// Save this newly created wallet address as the default address to
	// load up when the application finishes loading.
	walletAddress := strings.ToLower(account.Address.String())
	preferences.SetDefaultWalletAddress(strings.ToLower(walletAddress))

	// Return our address.
	return walletAddress, nil
}

func (a *App) SetDefaultWalletAddress(walletAddress string) {
	preferences := PreferencesInstance()
	preferences.SetDefaultWalletAddress(strings.ToLower(walletAddress))
}