package handler

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TokenTransferArgs struct {
	ChainID               uint16
	FromAccountAddress    *common.Address
	AccountWalletPassword string
	To                    *common.Address
	TokenID               *big.Int
}

type TokenTransferReply struct {
}

func (impl *ComicCoinRPCServer) TokenTransfer(args *TokenTransferArgs, reply *TokenTransferReply) error {
	err := impl.tokenTransferService.Execute(
		context.Background(),
		args.ChainID,
		args.FromAccountAddress,
		args.AccountWalletPassword,
		args.To,
		args.TokenID,
	)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = TokenTransferReply{}
	return nil
}
