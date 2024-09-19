package http

import (
	"log"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type httpInputPort struct {
	cfg        *config.Config
	kvs        keyvaluestore.KeyValueStorer
	blockchain *blockchain.Blockchain
}

func NewInputPort(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	bc *blockchain.Blockchain,
) inputport.InputPortServer {
	return &httpInputPort{
		cfg:        cfg,
		kvs:        kvs,
		blockchain: bc,
	}
}

func (s *httpInputPort) Run() {
	log.Println("TODO: Impl. Run()")
}

func (s *httpInputPort) Shutdown() {
	log.Println("TODO: Impl. Shutdown()")
}
