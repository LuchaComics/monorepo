package handler

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
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
	pass, secureErr := sstring.NewSecureString(args.AccountWalletPassword)
	if secureErr != nil {
		return secureErr
	}
	err := impl.coinTransferService.Execute(
		context.Background(),
		args.ChainID,
		args.FromAccountAddress,
		pass,
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
