package domain

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	TransactionTypeCoin  = "coin"
	TransactionTypeToken = "token"
)

// Transaction structure represents a transfer of coins between accounts
// which have not been added to the blockchain yet and are waiting for the miner
// to receive and verify. Once  transactions have been veriried
// they will be deleted from our system as they will live in the blockchain
// afterwords.
type Transaction struct {
	ChainID          uint16          `json:"chain_id"`           // Ethereum: The chain id that is listed in the genesis file.
	Nonce            uint64          `json:"nonce"`              // Ethereum: Unique id for the transaction supplied by the user.
	From             *common.Address `json:"from"`               // Ethereum: Account sending the transaction. Will be checked against signature.
	To               *common.Address `json:"to"`                 // Ethereum: Account receiving the benefit of the transaction.
	Value            uint64          `json:"value"`              // Ethereum: Monetary value received from this transaction.
	Tip              uint64          `json:"tip"`                // Ethereum: Tip offered by the sender as an incentive to mine this transaction.
	Data             []byte          `json:"data"`               // Ethereum: Extra data related to the transaction.
	Type             string          `json:"type"`               // ComicCoin: The type of transaction this is, either `coin` or `token`.
	TokenID          uint64          `json:"token_id"`           // ComicCoin: Unique identifier for the Token (if this transaciton is an Token).
	TokenMetadataURI string          `json:"token_metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
	TokenNonce       uint64          `json:"token_nonce"`        // ComicCoin: For every transaction action (mint, transfer, burn, etc), increment token nonce by value of 1.
}

// Sign function signs the  transaction using the user's private key
// and returns a signed version of that transaction.
func (tx Transaction) Sign(privateKey *ecdsa.PrivateKey) (SignedTransaction, error) {

	// Create a hash of the transaction to prepare it for signing.
	tran, err := tx.HashWithComicCoinStamp()
	if err != nil {
		return SignedTransaction{}, err
	}

	// Sign the hash using the private key to generate a signature.
	sig, err := crypto.Sign(tran, privateKey)
	if err != nil {
		return SignedTransaction{}, err
	}

	// Break the signature into the 3 parts: R, S, and V.
	v, r, s := toSignatureValues(sig)

	// Create the signed transaction, including the original transaction and the signature parts.
	signedTx := SignedTransaction{
		Transaction: tx,
		V:           v,
		R:           r,
		S:           s,
	}

	return signedTx, nil
}

// HashWithComicCoinStamp creates a unique hash of the transaction and
// prepares it for signing by adding a special "stamp".
func (tx Transaction) HashWithComicCoinStamp() ([]byte, error) {
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
func (tx SignedTransaction) FromAddress() (string, error) {

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
