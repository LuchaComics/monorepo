package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type ListBlockTransactionsByAddressArgs struct {
	AccountAddress *common.Address
}

type ListBlockTransactionsByAddressReply struct {
	BlockTransactions []*domain.BlockTransaction
}

func (impl *ComicCoinRPCServer) ListBlockTransactionsByAddress(args *ListBlockTransactionsByAddressArgs, reply *ListBlockTransactionsByAddressReply) error {
	txs, err := impl.listBlockTransactionsByAddressService.Execute(context.Background(), args.AccountAddress)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = ListBlockTransactionsByAddressReply{
		BlockTransactions: txs,
	}
	return nil
}
