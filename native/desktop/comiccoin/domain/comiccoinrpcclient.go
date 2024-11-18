package domain

import (
	"context"
	"math/big"
)

type ComicCoincRPCClient struct {
	// Empty
}

type ComicCoincRPCClientRepository interface {
	GetTimestamp(ctx context.Context) (uint64, error)
	GetNonFungibleToken(ctx context.Context, nftID *big.Int, directoryPath string) (*NonFungibleToken, error)
}
