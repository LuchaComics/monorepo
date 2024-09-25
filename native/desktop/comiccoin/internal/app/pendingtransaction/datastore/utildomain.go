package datastore

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// Sign function signs the pending transaction using the user's private key
// and returns a signed version of that transaction.
func (tx PendingTransaction) Sign(privateKey *ecdsa.PrivateKey) (SignedPendingTransaction, error) {

	// Create a hash of the transaction to prepare it for signing.
	tran, err := tx.HashWithComicCoinStamp()
	if err != nil {
		return SignedPendingTransaction{}, err
	}

	// Sign the hash using the private key to generate a signature.
	sig, err := crypto.Sign(tran, privateKey)
	if err != nil {
		return SignedPendingTransaction{}, err
	}

	// Break the signature into the 3 parts: R, S, and V.
	v, r, s := toSignatureValues(sig)

	// Create the signed transaction, including the original transaction and the signature parts.
	signedTx := SignedPendingTransaction{
		PendingTransaction: tx,
		V:                  v,
		R:                  r,
		S:                  s,
	}

	return signedTx, nil
}

// HashWithComicCoinStamp creates a unique hash of the transaction and
// prepares it for signing by adding a special "stamp".
func (tx PendingTransaction) HashWithComicCoinStamp() ([]byte, error) {
	// Convert the transaction into JSON format.
	txData, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	// Create a hash from the JSON transaction data.
	txHash := crypto.Keccak256Hash(txData)

	// Add a special stamp to identify this as a Comic Coin transaction.
	stamp := []byte("\x19Comic Coin Signed Message:\n32")
	tran := crypto.Keccak256Hash(stamp, txHash.Bytes())

	return tran.Bytes(), nil
}

const comicCoinID = 29

// toSignatureValues breaks down the 65-byte signature into its three components:
// R, S, and V. This format is used to help identify the blockchain that created the signature.
func toSignatureValues(sig []byte) (r, s, v *big.Int) {
	// Divide the signature into three parts: R, S, and V.
	// We add a unique number (comicCoinID) to V to identify the Comic Coin blockchain.
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + comicCoinID})

	return r, s, v
}

// Converts R, S, and V back into a single signature byte array.
func toSignatureBytes(v, r, s *big.Int) []byte {
	sig := make([]byte, crypto.SignatureLength)

	// Copy R, S into the signature array, and adjust V by subtracting the comicCoinID.
	copy(sig, r.Bytes())
	copy(sig[32:], s.Bytes())
	sig[64] = byte(v.Uint64() - comicCoinID)

	return sig
}

// Converts R, S, and V back into a signature, but leaves V unmodified.
func toSignatureBytesForDisplay(v, r, s *big.Int) []byte {
	sig := make([]byte, crypto.SignatureLength)

	copy(sig, r.Bytes())
	copy(sig[32:], s.Bytes())
	sig[64] = byte(v.Uint64())

	return sig
}

// FromAddress extracts the account address from the signed transaction by
// recovering the public key from the signature.
func (tx SignedPendingTransaction) FromAddress() (string, error) {

	// Create the hash of the transaction to prepare it for extracting the public key.
	tran, err := tx.HashWithComicCoinStamp()
	if err != nil {
		return "", err
	}

	// Combine R, S, and V into the original signature format.
	sig := toSignatureBytes(tx.V, tx.R, tx.S)

	// Use the signature to get the public key of the account that signed it.
	publicKey, err := crypto.SigToPub(tran, sig)
	if err != nil {
		return "", err
	}

	// Convert the public key to an account address.
	return crypto.PubkeyToAddress(*publicKey).String(), nil
}

// Validate checks if the transaction is valid. It verifies the signature,
// makes sure the account addresses are correct, and checks if the 'from'
// and 'to' accounts are not the same.
func (tx SignedPendingTransaction) Validate(chainID uint16) error {
	// Check if the transaction's chain ID matches the expected one.
	if tx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", tx.ChainID, chainID)
	}

	// Ensure the 'from' and 'to' accounts are not the same.
	if tx.From == tx.To {
		return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", tx.From, tx.To)
	}

	// Validate the signature parts (R, S, and V).
	if err := VerifySignature(tx.V, tx.R, tx.S); err != nil {
		return err
	}

	// Verify that the 'from' address matches the one from the signature.
	address, err := tx.FromAddress()
	if err != nil {
		return err
	}

	if address != string(tx.From.Hex()) {
		return errors.New("signature address doesn't match from address")
	}

	return nil
}

// VerifySignature checks if the signature is valid by ensuring the V value
// is correct and that the signature follows the proper rules.
func VerifySignature(v, r, s *big.Int) error {

	// Make sure V is either 0 or 1 (after subtracting the comicCoinID).
	uintV := v.Uint64() - comicCoinID
	if uintV != 0 && uintV != 1 {
		return errors.New("invalid recovery id")
	}

	// Check that R and S follow the signature rules.
	if !crypto.ValidateSignatureValues(byte(uintV), r, s, false) {
		return errors.New("invalid signature values")
	}

	return nil
}
