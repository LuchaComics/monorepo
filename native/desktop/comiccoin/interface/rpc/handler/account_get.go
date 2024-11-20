package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type GetAccountArgs struct {
	AccountAddress *common.Address
}

type GetAccountReply struct {
	Account *domain.Account
}

func (impl *ComicCoinRPCServer) GetAccount(args *GetAccountArgs, reply *GetAccountReply) error {

	account, err := impl.getAccountService.Execute(context.Background(), args.AccountAddress)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = GetAccountReply{
		Account: account,
	}
	return nil
}
