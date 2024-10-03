package domain

import "context"

type BlockchainIndexDTO string

type BlockchainIndexDTORepository interface {
	BroadcastRequestToP2PNetwork(ctx context.Context, dto *BlockchainIndexDTO) error
	ReceiveRequestFromP2PNetwork(ctx context.Context) (*BlockchainIndexDTO, error)
	BroadcastResponseToP2PNetwork(ctx context.Context, dto *BlockchainIndexDTO) error
	ReceiveResponseFromP2PNetwork(ctx context.Context) (*BlockchainIndexDTO, error)
}
