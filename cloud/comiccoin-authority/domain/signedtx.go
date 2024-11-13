package domain

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/fxamacker/cbor/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// SignedTransaction is a signed version of the transaction. This is how
// clients like a wallet provide transactions for inclusion into the blockchain.
type SignedTransaction struct {
	Transaction
	V *big.Int `bson:"v" json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with comicCoinID.
	R *big.Int `bson:"r" json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `bson:"s" json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// Validate checks if the transaction is valid. It verifies the signature,
// makes sure the account addresses are correct, and checks if the 'from'
// and 'to' accounts are not the same (unless you are the proof of authority!)
func (stx SignedTransaction) Validate(chainID uint16, isPoA bool) error {
	// Check if the transaction's chain ID matches the expected one.
	if stx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", stx.ChainID, chainID)
	}

	// Ensure the 'from' and 'to' accounts are not the same.
	if stx.From == stx.To {
		// ... unless you are the proof of authority.
		if !isPoA {
			return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", stx.From, stx.To)
		}
	}

	// Validate the signature parts (R, S, and V).
	if err := VerifySignature(stx.V, stx.R, stx.S); err != nil {
		return err
	}

	// Verify that the 'from' address matches the one from the signature.
	address, err := stx.FromAddress()
	if err != nil {
		return err
	}

	if address != string(stx.From.Hex()) {
		log.Printf("SignedTransaction: Validate: signature address %v doesn't match from address %v\n", address, stx.From.Hex())
		return errors.New("signature address doesn't match from address")
	}

	return nil
}

func (stx *SignedTransaction) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(stx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed transaction: %v", err)
	}
	return dataBytes, nil
}

func NewSignedTransactionFromDeserialize(data []byte) (*SignedTransaction, error) {
	// Variable we will use to return.
	stx := &SignedTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &stx); err != nil {
		return nil, fmt.Errorf("failed to deserialize signed transaction: %v", err)
	}
	return stx, nil
}

// MarshalBSON overrides the default serializer to handle a bug with mongodb.
func (stx *SignedTransaction) MarshalBSON() ([]byte, error) {
	// Developers note:
	// The reason *big.Int fields (like V, R, and S in SignedTransaction) aren't
	// being saved in MongoDB is because MongoDB's bson package does not natively
	// support encoding or decoding *big.Int values. By default, MongoDB doesn't
	// know how to handle big.Int types, so they end up being ignored.
	//
	// To fix this, you need to manually convert *big.Int values to a format
	// MongoDB can store (such as a string or integer) and then convert them back
	// on retrieval. Here’s one way to achieve this by adding custom serialization
	// for these fields

	type Alias SignedTransaction // Alias to avoid recursion
	return bson.Marshal(&struct {
		V string `bson:"v"`
		R string `bson:"r"`
		S string `bson:"s"`
		*Alias
	}{
		V:     stx.V.String(),
		R:     stx.R.String(),
		S:     stx.S.String(),
		Alias: (*Alias)(stx),
	})
}

// UnmarshalBSON overrides the default deserializer to handle a bug with mongodb.
func (stx *SignedTransaction) UnmarshalBSON(data []byte) error {
	// Developers note:
	// The reason *big.Int fields (like V, R, and S in SignedTransaction) aren't
	// being saved in MongoDB is because MongoDB's bson package does not natively
	// support encoding or decoding *big.Int values. By default, MongoDB doesn't
	// know how to handle big.Int types, so they end up being ignored.
	//
	// To fix this, you need to manually convert *big.Int values to a format
	// MongoDB can store (such as a string or integer) and then convert them back
	// on retrieval. Here’s one way to achieve this by adding custom serialization
	// for these fields

	type Alias SignedTransaction // Alias to avoid recursion
	aux := &struct {
		V string `bson:"v"`
		R string `bson:"r"`
		S string `bson:"s"`
		*Alias
	}{
		Alias: (*Alias)(stx),
	}

	if err := bson.Unmarshal(data, aux); err != nil {
		return err
	}

	stx.V, _ = new(big.Int).SetString(aux.V, 10)
	stx.R, _ = new(big.Int).SetString(aux.R, 10)
	stx.S, _ = new(big.Int).SetString(aux.S, 10)

	return nil
}
