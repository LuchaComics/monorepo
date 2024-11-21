package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetByBlockTransactionTimestampArgs struct {
	Timestamp uint64
}

type GetByBlockTransactionTimestampReply struct {
	BlockData *domain.BlockData
}

func (impl *ComicCoinRPCServer) GetByBlockTransactionTimestamp(args *GetByBlockTransactionTimestampArgs, reply *GetByBlockTransactionTimestampReply) error {

	blockData, err := impl.getByBlockTransactionTimestampService.Execute(context.Background(), args.Timestamp)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = GetByBlockTransactionTimestampReply{
		BlockData: blockData,
	}
	return nil
}
