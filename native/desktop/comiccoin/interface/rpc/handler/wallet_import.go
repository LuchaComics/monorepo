package handler

import (
	"context"
)

type ImportWalletArgs struct {
	FilePath string
}

type ImportWalletReply struct {
}

func (impl *ComicCoinRPCServer) ImportWallet(args *ImportWalletArgs, reply *ImportWalletReply) error {

	err := impl.importWalletService.Execute(context.Background(), args.FilePath)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = ImportWalletReply{}
	return nil
}
