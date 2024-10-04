package domain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

// MempoolTransactionDTO is the data-transfer object used by nodes to take
// the transactions submitted by their users and send and receive
// signed transaction via the peer-to-peer (P2P) network.
type MempoolTransactionDTO struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

type MempoolTransactionDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *MempoolTransactionDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*MempoolTransactionDTO, error)
}

func (dto *MempoolTransactionDTO) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize mempool transaction dto: %v", err)
	}
	return dataBytes, nil
}

func NewMempoolTransactionDTOFromDeserialize(data []byte) (*MempoolTransactionDTO, error) {
	// Variable we will use to return.
	dto := &MempoolTransactionDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("failed to deserialize mempool transaction dto: %v", err)
	}
	return dto, nil
}
