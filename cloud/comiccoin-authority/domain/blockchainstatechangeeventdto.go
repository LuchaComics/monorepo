package domain

import (
	"context"
)

// BlockchainStateChangeEventRepository is an interface that defines the methods for
// handling the blockchain state via the network using server sent events (SSE).
type BlockchainStateChangeEventDTORepository interface {
	SubscribeToBlockchainAuthority(ctx context.Context, chainID uint16) (<-chan uint16, error)
}
