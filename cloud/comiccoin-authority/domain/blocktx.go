package domain

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/fxamacker/cbor/v2"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
)

// BlockTransaction represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTransaction struct {
	SignedTransaction
	TimeStamp uint64 `bson:"timestamp" json:"timestamp"` // Ethereum: The time the transaction was received.
	Fee       uint64 `bson:"fee" json:"fee"`             // ComicCoin: Fee paid for this transaction to the ComicCoin authority.
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
	// Note: MongoDB doesn't support `*big.Int` so we are forced to do this.
	txV, txR, txS := tx.SignedTransaction.GetBigIntFields()
	otherTxV, otherTxR, otherTxS := otherTx.SignedTransaction.GetBigIntFields()

	txSig := signature.ToSignatureBytes(txV, txR, txS)
	otherTxSig := signature.ToSignatureBytes(otherTxV, otherTxR, otherTxS)

	return tx.GetNonce().Cmp(otherTx.GetNonce()) == 0 && bytes.Equal(txSig, otherTxSig)
}
