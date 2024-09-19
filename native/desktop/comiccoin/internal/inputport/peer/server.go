package peer

import (
	"log"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type peerInputPort struct {
	cfg        *config.Config
	kvs        keyvaluestore.KeyValueStorer
	blockchain *blockchain.Blockchain
	// host       host.Host
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

	return &peerInputPort{
		cfg:        cfg,
		kvs:        kvs,
		blockchain: bc,
		// host:       h,
	}
}

func (s *peerInputPort) Run() {
	log.Println("running peer")
	time.Sleep(10 * time.Second)

}

func (s *peerInputPort) Shutdown() {
	log.Println("shutting down")

}

// func makeHost(port int) (host.Host, error) {
// 	// Creates a new RSA key pair for this host.
// 	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
//
// 	// 0.0.0.0 will listen on any interface device.
// 	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
//
// 	// libp2p.New constructs a new libp2p Host.
// 	// Other options can be added here.
// 	return libp2p.New(
// 		libp2p.ListenAddrs(sourceMultiAddr),
// 		libp2p.Identity(prvKey),
// 	)
// }
