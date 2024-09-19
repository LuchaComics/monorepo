package ipfs

import (
	"log"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type ipfsInputPort struct {
	cfg        *config.Config
	kvs        keyvaluestore.KeyValueStorer
	blockchain *blockchain.Blockchain
}

func NewInputPort(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	bc *blockchain.Blockchain,
) inputport.InputPortServer {
	return &ipfsInputPort{
		cfg:        cfg,
		kvs:        kvs,
		blockchain: bc,
	}
}

func (s *ipfsInputPort) Run() {
	log.Println("running ipfs")
	time.Sleep(10 * time.Second)

}

func (s *ipfsInputPort) Shutdown() {

}
