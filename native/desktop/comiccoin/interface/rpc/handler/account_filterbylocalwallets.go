package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type AccountListingByLocalWalletsArgs struct {
}

type AccountListingByLocalWalletsReply struct {
	Accounts []*domain.Account
}

func (impl *ComicCoinRPCServer) AccountListingByLocalWallets(args *AccountListingByLocalWalletsArgs, reply *AccountListingByLocalWalletsReply) error {

	accounts, err := impl.accountListingByLocalWalletsService.Execute(context.Background())
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = AccountListingByLocalWalletsReply{
		Accounts: accounts,
	}
	return nil
}
