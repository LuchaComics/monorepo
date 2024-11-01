package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/ethereum/go-ethereum/common"
)

func (a *App) TransferToken(
	toRecipientAddress string,
	tokenID uint64,
	message string,
	senderAccountAddress string,
	senderAccountPassword string,
) error {
	// ctx := context.Background()

	a.logger.Debug("Transfering token...",
		slog.Any("toRecipientAddress", toRecipientAddress),
		slog.Any("tokenID", tokenID),
		slog.Any("message", message),
		slog.Any("senderAccountAddress", senderAccountAddress),
		slog.Any("senderAccountPassword", senderAccountPassword),
	)

	// var toRecipientAddr *common.Address = nil
	// if toRecipientAddress != "" {
	// 	to := common.HexToAddress(toRecipientAddress)
	// 	toRecipientAddr = &to
	// }
	//
	// var senderAccountAddr *common.Address = nil
	// if senderAccountAddress != "" {
	// 	sender := common.HexToAddress(senderAccountAddress)
	// 	senderAccountAddr = &sender
	// }
	//
	// err := a.transferCoinService.Execute(ctx, senderAccountAddr, senderAccountPassword, toRecipientAddr, coins, []byte(message))
	// if err != nil {
	// 	a.logger.Error("Failed transfering coin(s)",
	// 		slog.Any("error", err))
	// 	return err
	// }

	return nil
}

func (a *App) GetNonFungibleTokensByOwnerAddress(address string) ([]*domain.NonFungibleToken, error) {
	addr := common.HexToAddress(strings.ToLower(address))

	// Defensive code
	if address == "" {
		return make([]*domain.NonFungibleToken, 0), fmt.Errorf("failed because: address is null: %v", address)
	}

	//
	// STEP 1:
	// Lookup all the tokens. Note: A token only has `token_id` and
	// `metadata_uri` fields - nothing else!
	//

	toks, err := a.listNonFungibleTokensByOwnerService.Execute(&addr)
	if err != nil {
		a.logger.Error("Failed listing tokens by owner",
			slog.Any("error", err))
		return make([]*domain.NonFungibleToken, 0), err
	}

	a.logger.Debug("",
		slog.Any("toks", toks))

	return toks, nil
}
