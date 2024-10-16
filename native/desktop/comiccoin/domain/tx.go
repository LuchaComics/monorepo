package domain

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
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
	// Break the signature into the 3 parts: R, S, and V.
	v, r, s, err := signature.Sign(tx, privateKey)
	if err != nil {
		return SignedTransaction{}, err
	}

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
	return signature.HashWithComicCoinStamp(tx)
}

// FromAddress extracts the account address from the signed transaction by
// recovering the public key from the signature.
func (tx SignedTransaction) FromAddress() (string, error) {
	return signature.FromAddress(tx.Transaction, tx.V, tx.R, tx.S)
}

// VerifySignature checks if the signature is valid by ensuring the V value
// is correct and that the signature follows the proper rules.
func VerifySignature(v, r, s *big.Int) error {
	return signature.VerifySignature(v, r, s)
}
