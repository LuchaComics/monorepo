package main

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
)

func (a *App) TransferCoin(
	toRecipientAddress string,
	coins uint64,
	message string,
	senderAccountAddress string,
	senderAccountPassword string,
) error {
	ctx := context.Background()

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

	err := a.transferCoinService.Execute(ctx, senderAccountAddr, senderAccountPassword, toRecipientAddr, coins, []byte(message))
	if err != nil {
		a.logger.Error("Failed transfering coin(s)",
			slog.Any("error", err))
		return err
	}

	return nil
}
