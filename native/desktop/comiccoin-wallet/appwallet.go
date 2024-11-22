package main

import (
	"fmt"
	"log/slog"
	"strings"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin-wallet/common/preferences"
)

func (a *App) DefaultWalletAddress() string {
	preferences := pref.PreferencesInstance()
	return preferences.DefaultWalletAddress
}

func (a *App) ListWallets() ([]*domain.Wallet, error) {
	return a.walletsFilterByLocalService.Execute(a.ctx)
}

func (a *App) CreateWallet(walletPassword, walletPasswordRepeated, walletLabel string) (string, error) {
	pass, err := sstring.NewSecureString(walletPassword)
	if err != nil {
		a.logger.Error("Failed securing password",
			slog.Any("error", err))
		return "", err
	}
	defer pass.Wipe()
	passRepeated, err := sstring.NewSecureString(walletPasswordRepeated)
	if err != nil {
		a.logger.Error("Failed securing password repeated",
			slog.Any("error", err))
		return "", err
	}
	defer passRepeated.Wipe()

	account, err := a.createAccountService.Execute(a.ctx, pass, passRepeated, walletLabel)
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
	preferences := pref.PreferencesInstance()
	preferences.SetDefaultWalletAddress(strings.ToLower(walletAddress))
}