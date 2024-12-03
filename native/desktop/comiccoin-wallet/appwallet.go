package main

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) DefaultWalletAddress() string {
	preferences := PreferencesInstance()
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
	// defer pass.Wipe() // Developers Note: Commented out b/c they are causing problems with our app.
	passRepeated, err := sstring.NewSecureString(walletPasswordRepeated)
	if err != nil {
		a.logger.Error("Failed securing password repeated",
			slog.Any("error", err))
		return "", err
	}
	// defer passRepeated.Wipe() // Developers Note: Commented out b/c they are causing problems with our app.

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
	preferences := PreferencesInstance()
	preferences.SetDefaultWalletAddress(strings.ToLower(walletAddress))
}

func (a *App) ExportWalletUsingDialog(walletAddress string) error {
	// Initialize Wails runtime
	result, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Please select were to save the blockchain",
	})
	if err != nil {
		a.logger.Error("Failed opening directory dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	a.logger.Debug("->", slog.Any("result", result))
	return nil
}
