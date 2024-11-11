package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/peer"
)

// BlockDataDTO is the data-transfer object used by nodes to send back and forth
// BlockData across the distributed / P2P network for the blockchain.
// It contains the hash of the block, the block header, and the list of transactions in the block.
type BlockDataDTO struct {
	// Hash is the unique hash of the block.
	Hash string `json:"hash"`

	// The signature of this block's "Header" field which was applied by the
	// proof-of-authority validator.
	HeaderSignatureBytes []byte `json:"header_signature_bytes"`

	// Header is the block header, which contains metadata about the block.
	Header *BlockHeader `json:"block_header"`

	// Trans is the list of transactions in the block.
	Trans []BlockTransaction `json:"trans"`

	// The proof-of-authority validator whom executed the validation of
	// this block data in our blockchain.
	Validator *Validator `json:"validator"`
}

// BlockDataDTORepository is an interface that defines the methods for interacting with block data DTOs.
// It provides methods for sending requests to peers, receiving requests from peers, sending responses to peers, and receiving responses from peers.
type BlockDataDTORepository interface {
	// SendRequestToRandomPeer sends a request to a random connected peer in the peer-to-peer network.
	// It takes a context and the hash of the block data and returns an error if one occurs.
	SendRequestToRandomPeer(ctx context.Context, blockDataHash string) error

	// ReceiveRequestFromNetwork receives a request from the peer-to-peer network.
	// It takes a context and returns the peer ID of the peer that sent the request, the hash of the block data, and an error if one occurs.
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, string, error)

	// SendResponseToPeer sends a response to a peer that requested block data.
	// It takes a context, the peer ID of the peer to send the response to, and the block data DTO to send.
	// It returns an error if one occurs.
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data *BlockDataDTO) error

	// ReceiveResponseFromNetwork receives a response from the peer-to-peer network.
	// It takes a context and returns the block data DTO and an error if one occurs.
	ReceiveResponseFromNetwork(ctx context.Context) (*BlockDataDTO, error)
}

// Serialize serializes a block data DTO into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockDataDTO) Serialize() ([]byte, error) {
	// Marshal the block data DTO into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data dto: %v", err)
	}
	return dataBytes, nil
}

// NewBlockDataDTOFromDeserialize deserializes a block data DTO from a byte array.
// It returns the deserialized block data DTO and an error if one occurs.
func NewBlockDataDTOFromDeserialize(data []byte) (*BlockDataDTO, error) {
	// Variable we will use to return.
	blockDataDTO := &BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data DTO using CBOR.
	if err := cbor.Unmarshal(data, &blockDataDTO); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data dto: %v", err)
	}
	return blockDataDTO, nil
}
