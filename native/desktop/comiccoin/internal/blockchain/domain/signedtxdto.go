package domain

import (
	"context"
	"math/big"
)

// SignedTransactionDTO is the data-transfer object used to send and receive
// signed transaction via the peer-to-peer (P2P) network.
type SignedTransactionDTO struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

type SignedTransactionDTORepository interface {
	Broadcast(ctx context.Context, bd *SignedTransactionDTO) error
	Receive(ctx context.Context) (*SignedTransactionDTO, error)
}
