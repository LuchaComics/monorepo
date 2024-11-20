package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type CreateAccountArgs struct {
	Password         string
	PasswordRepeated string
	Label            string
}

type CreateAccountReply struct {
	Account *domain.Account
}

func (impl *ComicCoinRPCServer) CreateAccount(args *CreateAccountArgs, reply *CreateAccountReply) error {

	account, err := impl.createAccountService.Execute(context.Background(), args.Password, args.PasswordRepeated, args.Label)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = CreateAccountReply{
		Account: account,
	}
	return nil
}
