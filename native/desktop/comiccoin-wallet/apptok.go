package main

import (
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	comic_domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	// auth_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

func (a *App) TransferToken(
	toRecipientAddress string,
	tokenID *big.Int,
	senderAccountAddress string,
	senderAccountPassword string,
) error {

	a.logger.Debug("Transfering token...",
		slog.Any("toRecipientAddress", toRecipientAddress),
		slog.Any("tokenID", tokenID),
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

	password, err := sstring.NewSecureString(senderAccountPassword)
	if err != nil {
		a.logger.Error("Failed securing password",
			slog.Any("error", err))
		return err
	}
	defer password.Wipe()

	tokenTransferErr := a.tokenTransferService.Execute(
		a.ctx,
		preferences.ChainID,
		senderAccountAddr,
		password,
		toRecipientAddr,
		tokenID,
	)
	if tokenTransferErr != nil {
		a.logger.Error("Failed transfering token",
			slog.Any("error", tokenTransferErr))
		return tokenTransferErr
	}

	return nil
}

func (a *App) GetNonFungibleTokensByOwnerAddress(address string) ([]*comic_domain.NonFungibleToken, error) { //TODO FINISH BELOW
	// addr := common.HexToAddress(strings.ToLower(address))
	//
	// // Defensive code
	// if address == "" {
	// 	return make([]*comic_domain.NonFungibleToken, 0), fmt.Errorf("failed because: address is null: %v", address)
	// }
	//
	// //
	// // STEP 1:
	// // Lookup all the tokens. Note: A token only has `token_id` and
	// // `metadata_uri` fields - nothing else!
	// //
	//
	// toks, err := a.tokenListByOwnerService.Execute(a.ctx, &addr)
	// if err != nil {
	// 	a.logger.Error("Failed listing tokens by owner",
	// 		slog.Any("error", err))
	// 	return make([]*comic_domain.NonFungibleToken, 0), err
	// }
	//
	// a.logger.Debug("",
	// 	slog.Any("toks_count", len(toks)))
	//
	// return toks, nil
	return make([]*comic_domain.NonFungibleToken, 0), nil //TODO: REPAIR ABOVE.
}

func (a *App) GetNonFungibleToken(tokenID *big.Int) (*comic_domain.NonFungibleToken, error) { //TODO FINISH BELOW
	// nftok, err := a.getOrDownloadNonFungibleTokenService.Execute(a.ctx, tokenID)
	// if err != nil {
	// 	a.logger.Error("Failed getting non-fungible token by token ID.",
	// 		slog.Any("error", err))
	// 	return nil, err
	// }
	//
	// a.logger.Debug("",
	// 	slog.Any("nftok", nftok))
	//
	// return nftok, nil
	return nil, nil //TODO: REPAIR ABOVE.
}
