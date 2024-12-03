package main

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
)

func (a *App) TransferCoin(
	toRecipientAddress string,
	coins uint64,
	message string,
	senderAccountAddress string,
	senderAccountPassword string,
) error {
	//
	// STEP 1: Validation.
	//

	a.logger.Debug("Transfering coin...",
		slog.Any("toRecipientAddress", toRecipientAddress),
		slog.Any("coins", coins),
		slog.Any("message", message),
		slog.Any("senderAccountAddress", senderAccountAddress),
		slog.Any("senderAccountPassword", senderAccountPassword),
	)

	e := make(map[string]string)
	if toRecipientAddress == "" {
		e["from_account_address"] = "missing value"
	}
	if senderAccountPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if toRecipientAddress == "" {
		e["to"] = "missing value"
	}
	if coins == 0 {
		e["value"] = "missing value"
	}
	if len(e) != 0 {
		a.logger.Warn("Failed validating create transaction parameters",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	var toRecipientAddr *common.Address = nil
	if toRecipientAddress != "" {
		to := common.HexToAddress(toRecipientAddress)
		toRecipientAddr = &to
	}

	var senderAccountAddr *common.Address = nil
	if senderAccountAddress != "" {
		sender := common.HexToAddress(senderAccountAddress)
		senderAccountAddr = &sender
	}

	preferences := PreferencesInstance()

	password, err := sstring.NewSecureString(senderAccountPassword)
	if err != nil {
		e := make(map[string]string)
		e["senderAccountPassword"] = "missing value"
		a.logger.Error("Failed securing password",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}
	// defer password.Wipe()  Developers Note: Commented out b/c they are causing problems with our app.

	coinTransferErr := a.coinTransferService.Execute(
		a.ctx,
		preferences.ChainID,
		senderAccountAddr,
		password,
		toRecipientAddr,
		coins,
		[]byte(message),
	)
	if coinTransferErr != nil {
		a.logger.Error("Failed transfering coin(s)",
			slog.Any("error", coinTransferErr))
		return coinTransferErr
	}

	return nil
}
