package p2p

import (
	"log"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type nodeInputPort struct {
	cfg                  *config.Config
	logger               *slog.Logger
	keypairStorer        keypair_ds.KeypairStorer
	blockchainController blockchain_c.BlockchainController
	host                 host.Host
	peerAddresses        []*peer.AddrInfo
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	kp keypair_ds.KeypairStorer,
	bc blockchain_c.BlockchainController,
) inputport.InputPortServer {

	//
	// CASE 1 OF 2:
	// HOST MODE
	//

	if cfg.Peer.BootstrapPeers == "" {
		logger.Info("Starting p2p node", slog.String("mode", "host"))

		// DEVELOPERS NOTE:
		// A Host contains all the core functionalities you require connecting
		// one peer to another. A Host contains an ID which is the identity of
		// a peer. The Host also contains a Peerstore which is like an address
		// book. With a Host you can Connect to other peers and create Streams
		// between them. A Stream represents a communication channel between
		// two peers in a libp2p network.
		//
		// Link :https://ldej.nl/post/building-an-echo-application-with-libp2p/

		node := &nodeInputPort{
			cfg:                  cfg,
			logger:               logger,
			keypairStorer:        kp,
			blockchainController: bc,
		}

		host, err := node.newHostWithPredictableIdentifier()
		if err != nil {
			log.Fatal(err)
		}
		node.host = host

		return node
	}

	//
	// CASE 2 OF 2:
	// DIAL MODE
	//

	logger.Info("Starting p2p node", slog.String("mode", "dial"))

	return &nodeInputPort{
		cfg:                  cfg,
		logger:               logger,
		keypairStorer:        kp,
		blockchainController: bc,
	}
}

func (node *nodeInputPort) Run() {
	node.logger.Info("Running p2p node")

}

func (node *nodeInputPort) Shutdown() {
	node.logger.Info("Gracefully shutting down p2p node")
	node.host.Close()
}
