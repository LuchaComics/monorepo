package domain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

// BlockchainState represents the entire blockchain state at the current moment
// in time of operation.
type BlockchainState struct {
	// The unique identifier for this blockchain that we are managing the state for.
	ChainID uint16 `bson:"chain_id" json:"chain_id"`

	LatestBlockNumberBytes []byte `bson:"latest_block_number_bytes" json:"latest_block_number_bytes"`
	LatestHash             string `bson:"latest_hash" json:"latest_hash"`
	LatestTokenIDBytes     []byte `bson:"latest_token_id_bytes" json:"latest_token_id_bytes"`

	// The current transaction fee to be apply for every transaction as of this moment in time that is received by the Authority from the Global Blockchain Network.
	TransactionFee uint64 `bson:"transaction_fee" json:"transaction_fee"`

	AccountHashState string `bson:"account_hash_state" json:"account_hash_state"`
	TokenHashState   string `bson:"token_hash_state" json:"token_hash_state"`
}

func (bs *BlockchainState) GetLatestBlockNumber() *big.Int {
	return new(big.Int).SetBytes(bs.LatestBlockNumberBytes)
}

func (bs *BlockchainState) IsLatestBlockNumberZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(bs.GetLatestBlockNumber().Bits()) == 0
}

func (bs *BlockchainState) SetLatestBlockNumber(n *big.Int) {
	bs.LatestBlockNumberBytes = n.Bytes()
}

func (bs *BlockchainState) GetLatestTokenID() *big.Int {
	return new(big.Int).SetBytes(bs.LatestTokenIDBytes)
}

func (bs *BlockchainState) IsLatestTokenIDZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(bs.GetLatestTokenID().Bits()) == 0
}

func (bs *BlockchainState) SeLatestTokenID(n *big.Int) {
	bs.LatestTokenIDBytes = n.Bytes()
}

type BlockchainStateRepository interface {
	// Upsert inserts or updates an blockchain state in the repository.
	UpsertByChainID(ctx context.Context, acc *BlockchainState) error

	// GetByChainID retrieves an blockchain state by its chain ID.
	GetByChainID(ctx context.Context, chainID uint16) (*BlockchainState, error)

	// ListAll retrieves all blockchain states in the repository.
	ListAll(ctx context.Context) ([]*BlockchainState, error)

	// DeleteByChainID deletes an blockchain state by its chain ID.
	DeleteByChainID(ctx context.Context, chainID uint16) error

	GetUpdateChangeStreamChannel(ctx context.Context) (<-chan BlockchainState, chan struct{}, error)

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Serialize serializes the blockchain state into a byte slice.
// This method uses the cbor library to marshal the blockchainState into a byte slice.
func (b *BlockchainState) Serialize() ([]byte, error) {
	// Marshal the blockchainState into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize blockchainState: %v", err)
	}
	return dataBytes, nil
}

// NewBlockchainStateFromDeserialize deserializes an blockchainState from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an blockchainState.
func NewBlockchainStateFromDeserialize(data []byte) (*BlockchainState, error) {
	// Create a new blockchainState variable to return.
	blockchainState := &BlockchainState{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the blockchainState variable using the cbor library.
	if err := cbor.Unmarshal(data, &blockchainState); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize blockchainState: %v", err)
	}
	return blockchainState, nil
}
