package handler

import (
	"context"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetNonFungibleTokenArgs struct {
	NonFungibleTokenID *big.Int
	DirectoryPath      string
}

type GetNonFungibleTokenReply struct {
	NonFungibleToken *domain.NonFungibleToken
}

func (impl *ComicCoinRPCServer) GetNonFungibleToken(args *GetNonFungibleTokenArgs, reply *GetNonFungibleTokenReply) error {

	nft, err := impl.getOrDownloadNonFungibleTokenService.Execute(context.Background(), args.NonFungibleTokenID, args.DirectoryPath)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = GetNonFungibleTokenReply{
		NonFungibleToken: nft,
	}
	return nil
}
