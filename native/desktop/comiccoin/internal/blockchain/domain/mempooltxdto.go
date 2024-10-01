package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/big"
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
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewMempoolTransactionDTOFromDeserialize(data []byte) (*MempoolTransactionDTO, error) {
	// Variable we will use to return.
	dto := &MempoolTransactionDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return dto, nil
}
