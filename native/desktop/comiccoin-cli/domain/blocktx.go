package domain

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/fxamacker/cbor/v2"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/blockchain/signature"
)

// BlockTransaction represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTransaction struct {
	SignedTransaction
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}

func (dto *BlockTransaction) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block transaction: %v", err)
	}
	return dataBytes, nil
}

func NewBlockTransactionFromDeserialize(data []byte) (*BlockTransaction, error) {
	// Variable we will use to return.
	dto := &BlockTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("failed to deserialize block transaction: %v", err)
	}
	return dto, nil
}

// Hash implements the merkle Hashable interface for providing a hash
// of a block transaction.
func (tx BlockTransaction) Hash() ([]byte, error) {
	str := signature.Hash(tx)

	// Need to remove the 0x prefix from the hash.
	return hex.DecodeString(str[2:])
}

// Equals implements the merkle Hashable interface for providing an equality
// check between two block transactions. If the nonce and signatures are the
// same, the two blocks are the same.
func (tx BlockTransaction) Equals(otherTx BlockTransaction) bool {
	txSig := signature.ToSignatureBytes(tx.V, tx.R, tx.S)
	otherTxSig := signature.ToSignatureBytes(otherTx.V, otherTx.R, otherTx.S)

	return tx.Nonce == otherTx.Nonce && bytes.Equal(txSig, otherTxSig)
}
