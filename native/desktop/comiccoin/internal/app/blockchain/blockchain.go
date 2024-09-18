package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
)

type AuthorizedNode struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

func (node *AuthorizedNode) Sign(data []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, node.PrivateKey, data)
}

// // PoABlockchain struct represents a Blochcain based on a `proof of authority`
// // consensus model.
// type PoABlockchain struct {
// 	genesisBlock     *Block
// 	chain            []*Block
// 	difficulty       int
// 	AuthorizedNodes  []*AuthorizedNode
// 	CurrentNodeIndex int
// }
//
// func NewPoABlockchain(authorizedKey *keystore.Key) (*PoABlockchain, error) {
// 	blockchain := &PoABlockchain{
// 		chain:            []*Block{NewGenesisBlock()},
// 		AuthorizedNodes:  make([]*AuthorizedNode, 1),
// 		CurrentNodeIndex: 0,
// 	}
//
// 	blockchain.AuthorizedNodes[0] = &AuthorizedNode{
// 		Address:    authorizedKey.Address.Hex(),
// 		PrivateKey: authorizedKey.PrivateKey,
// 	}
//
// 	return blockchain, nil
// }

// PoABlockchain struct represents a Blochcain based on a `proof of authority`
// consensus model.
type PoABlockchain struct {
	authorizedKey *keystore.Key
	genesisBlock  *Block
	chain         []*Block
	difficulty    int
	Database      keyvaluestore.KeyValueStorer
	lastHash      string
}

func NewPoABlockchain(kvs keyvaluestore.KeyValueStorer, authorizedKey *keystore.Key) *PoABlockchain {
	//
	// STEP 1:
	// Look up our `Genesis` block in our database and if it does not exist then
	// we need to create it and save it to the database.
	//

	var genesisBlock *Block
	gensisBlockBin, err := kvs.Get([]byte("genesis"))
	if err != nil {
		log.Printf("failed getting `gensis` from kvs: %v", err)
	}
	fmt.Println("gensisBlockBin ->>", gensisBlockBin)

	if gensisBlockBin == nil {
		genesisBlock = &Block{
			Hash:      "0",
			Timestamp: time.Now(),
		}
		err := kvs.Set([]byte("genesis"), genesisBlock.Serialize())
		if err != nil {
			log.Printf("failed setting `gensis` to kvs: %v", err)
		}
		fmt.Println("created genesis")
	} else {
		genesisBlock = Deserialize(gensisBlockBin)
		fmt.Println("fetched genesis")
	}

	//
	// STEP 2:
	// Lookup our latest block hash and if it doesn't exist then we will set
	// our `Genesis` block as our latest block hash
	//

	// We need to always keep a record of the last Hash for our application.
	lastHashBin, err := kvs.Get([]byte("lh"))
	if err != nil {
		log.Printf("failed getting `lh` from kvs: %v", err)
	}
	lastHash := string(lastHashBin)
	fmt.Println("lastHash ->>", lastHash)

	//
	//
	//

	difficulty := 1

	//
	//
	//

	return &PoABlockchain{
		authorizedKey,
		genesisBlock,
		[]*Block{genesisBlock},
		difficulty,
		kvs,
		lastHash,
	}
}

func (b *PoABlockchain) AddBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	fmt.Println(newBlock.Hash)

	if err := b.Database.Set([]byte(fmt.Sprintf("hash_%v", newBlock.Hash)), newBlock.Serialize()); err != nil {
		log.Fatalf("failed to set db: %v", err)
	}
	if err := b.Database.Set([]byte("lh"), []byte(newBlock.Hash)); err != nil {
		log.Fatalf("failed to set db: %v", err)
	}

	b.chain = append(b.chain, &newBlock)
}

func (b PoABlockchain) IsValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// func (chain *PoABlockchain) AddBlock(Data string) {
// 	prevBlock := chain.blocks[len(chain.blocks)-1]
// 	new := CreateBlock(Data, prevBlock.Hash)
// 	chain.blocks = append(chain.blocks, new)
// }
//
// func (chain *PoABlockchain) GetBlocks() []*Block {
// 	return chain.blocks
// }
