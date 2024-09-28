package domain

import "math/big"

// MempoolTransaction is a signed version of the transaction that is pending
// proof of work and verification by the blockchain network. This record is
// is to wait in the Mempool until the miner processes it.
type MempoolTransaction struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

type MempoolTransactionRepository interface {
	Upsert(bd *MempoolTransaction) error
	GetByNonce(nonce uint64) (*MempoolTransaction, error)
	List() ([]*MempoolTransaction, error)
	DeleteByNonce(nonce uint64) error
}
