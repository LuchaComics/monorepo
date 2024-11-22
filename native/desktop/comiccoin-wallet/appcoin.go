package main

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
)

func (a *App) TransferCoin(
	toRecipientAddress string,
	coins uint64,
	message string,
	senderAccountAddress string,
	senderAccountPassword string,
) error {
	a.logger.Debug("Transfering coin...",
		slog.Any("toRecipientAddress", toRecipientAddress),
		slog.Any("coins", coins),
		slog.Any("message", message),
		slog.Any("senderAccountAddress", senderAccountAddress),
		slog.Any("senderAccountPassword", senderAccountPassword),
	)

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
		a.logger.Error("Failed securing password",
			slog.Any("error", err))
		return err
	}
	defer password.Wipe()

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
