package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type ListTokensByOwnerAddressArgs struct {
	OwnerAddress *common.Address
}

type ListTokensByOwnerAddressReply struct {
	Tokens []*domain.Token
}

func (impl *ComicCoinRPCServer) ListTokensByOwnerAddress(args *ListTokensByOwnerAddressArgs, reply *ListTokensByOwnerAddressReply) error {
	toks, err := impl.tokenListByOwnerService.Execute(context.Background(), args.OwnerAddress)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = ListTokensByOwnerAddressReply{
		Tokens: toks,
	}
	return nil
}
