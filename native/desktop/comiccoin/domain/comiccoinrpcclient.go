package domain

import (
	"context"
)

type ComicCoincRPCClient struct {
}

type ComicCoincRPCClientRepository interface {
	GetTimestamp(ctx context.Context) (uint64, error)
}
