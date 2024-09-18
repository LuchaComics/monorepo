package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type AuthorizedNode struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

func (node *AuthorizedNode) Sign(data []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, node.PrivateKey, data)
}

// PoABlockchain struct represents a Blochcain based on a `proof of authority`
// consensus model.
type PoABlockchain struct {
	genesisBlock     *Block
	chain            []*Block
	difficulty       int
	AuthorizedNodes  []*AuthorizedNode
	CurrentNodeIndex int
}

func NewPoABlockchain(authorizedKey *keystore.Key) (*PoABlockchain, error) {
	blockchain := &PoABlockchain{
		chain:            []*Block{NewGenesisBlock()},
		AuthorizedNodes:  make([]*AuthorizedNode, 1),
		CurrentNodeIndex: 0,
	}

	blockchain.AuthorizedNodes[0] = &AuthorizedNode{
		Address:    authorizedKey.Address.Hex(),
		PrivateKey: authorizedKey.PrivateKey,
	}

	return blockchain, nil
}

func (b *PoABlockchain) AddBlock(from, to string, amount float64) error {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := &Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
	return nil
}
func (b *PoABlockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}
