package handler

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type ExportWalletArgs struct {
	AccountAddress *common.Address
	FilePath       string
}

type ExportWalletReply struct {
}

func (impl *ComicCoinRPCServer) ExportWallet(args *ExportWalletArgs, reply *ExportWalletReply) error {

	err := impl.exportWalletService.Execute(context.Background(), args.AccountAddress, args.FilePath)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = ExportWalletReply{}
	return nil
}
