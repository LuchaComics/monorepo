package p2p

import (
	"context"
	"log"
	"log/slog"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"

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
	kademliaDHT          *dht.IpfsDHT
	routingDiscovery     *routing.RoutingDiscovery
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	kp keypair_ds.KeypairStorer,
	bc blockchain_c.BlockchainController,
) inputport.InputPortServer {
	ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
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

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	node.host.SetStreamHandler(fetchProtocolVersion, func(stream network.Stream) {
		node.logger.Info("Got a new stream!")
		go NewFetchProtocol(node, stream)
		// 'stream' will stay open until you close it (or the other side closes it).
	})

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	//
	// Source: https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/chat.go#L112
	kademliaDHT, err := node.newKademliaDHT(ctx)
	if err != nil {
		logger.Debug("Failed creating new kademlia dht",
			slog.Any("error", err))
		log.Fatal(err)
	}
	node.kademliaDHT = kademliaDHT

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	node.logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(node.kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, node.cfg.Peer.RendezvousString)
	node.logger.Debug("Successfully announced!")
	node.routingDiscovery = routingDiscovery

	return node
}

func (node *nodeInputPort) Run() {
	ctx := context.Background()
	node.logger.Info("Running p2p node")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	node.logger.Debug("Searching for other peers...")
	peerChan, err := node.routingDiscovery.FindPeers(ctx, node.cfg.Peer.RendezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == node.host.ID() {
			continue
		}
		node.logger.Debug("Found peer:",
			slog.Any("peer", peer))

		node.logger.Debug("Connecting to:",
			slog.Any("peer", peer))

		stream, err := node.host.NewStream(ctx, peer.ID, fetchProtocolVersion)
		if err != nil {
			node.logger.Warn("Connection failed:",
				slog.Any("error", err))
			continue
		} else {
			node.logger.Info("Got a new stream!")
			go NewFetchProtocol(node, stream)
			// 'stream' will stay open until you close it (or the other side closes it).
		}

		node.logger.Info("Connected to:",
			slog.Any("peer", peer))
	}

}

func (node *nodeInputPort) Shutdown() {
	node.logger.Info("Gracefully shutting down p2p node")
	node.host.Close()
}

func (node *nodeInputPort) handleStream(stream network.Stream) {
	node.logger.Info("Got a new stream!")
	worker := NewServiceWorker(node, stream)
	go worker.Start(context.Background())
	// 'stream' will stay open until you close it (or the other side closes it).
}
