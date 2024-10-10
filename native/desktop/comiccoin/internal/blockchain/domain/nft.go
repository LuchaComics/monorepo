package domain

import (
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
	"github.com/ethereum/go-ethereum/common"
)

// NFTTransaction represents a transaction specifically related to an NFT.
type NFTTransaction struct {
	ChainID   uint16          `json:"chain_id"`  // Ethereum: The chain id that is listed in the genesis file.
	TokenID   uint64          `json:"token_id"`  // Unique identifier for the NFT.
	From      *common.Address `json:"from"`      // Account sending the NFT.
	To        *common.Address `json:"to"`        // Account receiving the NFT.
	Metadata  string          `json:"metadata"`  // URI pointing to NFT metadata file.
	TimeStamp uint64          `json:"timestamp"` // Timestamp of the NFT transaction.
}

// SignedNFTTransaction is a signed version of the transaction. This is how
// clients like a wallet provide transactions for inclusion into the blockchain.
type SignedNFTTransaction struct {
	NFTTransaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// Sign function signs the  transaction using the user's private key
// and returns a signed version of that transaction.
func (tx NFTTransaction) Sign(privateKey *ecdsa.PrivateKey) (SignedNFTTransaction, error) {

	// Create a hash of the transaction to prepare it for signing.
	txHash, err := signature.HashWithComicCoinStamp(tx)
	if err != nil {
		log.Printf("NFTTransaction: SIGN: signature.HashWithComicCoinStamp err: %v\n", err)
		return SignedNFTTransaction{}, err
	}

	// Sign the hash using the private key to generate a signature.
	v, r, s, err := signature.Sign(txHash, privateKey)
	if err != nil {
		log.Printf("NFTTransaction: SIGN: signature.Sign err: %v\n", err)
		return SignedNFTTransaction{}, err
	}

	// Create the signed transaction, including the original transaction and the signature parts.
	signedTx := SignedNFTTransaction{
		NFTTransaction: tx,
		V:              v,
		R:              r,
		S:              s,
	}

	return signedTx, nil
}
