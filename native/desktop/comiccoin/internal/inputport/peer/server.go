package peer

import (
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type peerInputPort struct {
	cfg        *config.Config
	kvs        keyvaluestore.KeyValueStorer
	blockchain *blockchain.Blockchain
	host       host.Host
}

func NewInputPort(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	bc *blockchain.Blockchain,
) inputport.InputPortServer {

	// // Create a new P2P host on port 9000
	// h, err := makeHost(cfg.AppPort)
	// if err != nil {
	// 	log.Fatalf("failed making this node a p2p host: %v", err)
	// }

	sourcePort := 9000

	h, err := makeHost(sourcePort)
	if err != nil {
		log.Fatalf("failed to make p2p host: %v", err)
	}

	return &peerInputPort{
		cfg:        cfg,
		kvs:        kvs,
		blockchain: bc,
		host:       h,
	}
}

func (s *peerInputPort) Run() {
	log.Println("running peer")
	time.Sleep(10 * time.Second)

}

func (s *peerInputPort) Shutdown() {
	log.Println("shutting down")
	s.host.Close()
}

func makeHost(port int) (host.Host, error) {
	node, err := libp2p.New()
	if err != nil {
		return nil, err
	}

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())

	return node, nil
}
