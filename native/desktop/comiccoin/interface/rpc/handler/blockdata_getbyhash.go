package handler

import (
	"context"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockDataGetByHashArgs struct {
	Hash string
}

type BlockDataGetByHashReply struct {
	BlockData *domain.BlockData
}

func (impl *ComicCoinRPCServer) BlockDataGetByHash(args *BlockDataGetByHashArgs, reply *BlockDataGetByHashReply) error {

	bd, err := impl.blockDataGetByHashService.Execute(context.Background(), args.Hash)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = BlockDataGetByHashReply{
		BlockData: bd,
	}
	return nil
}
