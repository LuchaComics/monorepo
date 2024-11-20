package handler

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type CoinTransferArgs struct {
	ChainID               uint16
	FromAccountAddress    *common.Address
	AccountWalletPassword string
	To                    *common.Address
	Value                 uint64
	Data                  []byte
}

type CoinTransferReply struct {
}

func (impl *ComicCoinRPCServer) CoinTransfer(args *CoinTransferArgs, reply *CoinTransferReply) error {
	err := impl.coinTransferService.Execute(
		context.Background(),
		args.ChainID,
		args.FromAccountAddress,
		args.AccountWalletPassword,
		args.To,
		args.Value,
		args.Data,
	)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = CoinTransferReply{}
	return nil
}
