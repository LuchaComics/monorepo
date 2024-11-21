package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
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
	pass, err := sstring.NewSecureString(args.Password)
	if err != nil {
		return err
	}
	passRepeated, err := sstring.NewSecureString(args.PasswordRepeated)
	if err != nil {
		return err
	}

	account, err := impl.createAccountService.Execute(context.Background(), pass, passRepeated, args.Label)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = CreateAccountReply{
		Account: account,
	}
	return nil
}
