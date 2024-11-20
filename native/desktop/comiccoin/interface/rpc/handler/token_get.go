package handler

import (
	"context"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetTokenArgs struct {
	TokenID *big.Int
}

type GetTokenReply struct {
	Token *domain.Token
}

func (impl *ComicCoinRPCServer) GetToken(args *GetTokenArgs, reply *GetTokenReply) error {

	token, err := impl.tokenGetService.Execute(context.Background(), args.TokenID)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = GetTokenReply{
		Token: token,
	}
	return nil
}
