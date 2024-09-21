package p2p

import (
	"bufio"
	"context"
	"log"
	"log/slog"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
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
	peerAddresses        []*peer.AddrInfo
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	kp keypair_ds.KeypairStorer,
	bc blockchain_c.BlockchainController,
) inputport.InputPortServer {

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

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	node.host.SetStreamHandler("/p2p/1.0.0", node.handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	var options []dht.Option
	bootstrapPeers := make([]peer.AddrInfo, len(node.cfg.Peer.BootstrapPeers))
	for i, addr := range node.cfg.Peer.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		bootstrapPeers[i] = *peerinfo
	}
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
		node.logger.Info("Running p2p node as dht server")
	}
	options = append(options, dht.BootstrapPeers(bootstrapPeers...))

	kademliaDHT, err := dht.New(ctx, host, options...)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	logger.Debug("Bootstrapping the DHT",
		slog.String("RendezvousString", node.cfg.Peer.RendezvousString))
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Wait a bit to let bootstrapping finish (really bootstrap should block until it's ready, but that isn't the case yet.)
	time.Sleep(1 * time.Second)

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, node.cfg.Peer.RendezvousString)
	logger.Debug("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	logger.Debug("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, node.cfg.Peer.RendezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		logger.Debug("Found peer:",
			slog.Any("peer", peer))

		logger.Debug("Connecting to:",
			slog.Any("peer", peer))

		stream, err := host.NewStream(ctx, peer.ID, "/p2p/1.0.0")

		if err != nil {
			logger.Warn("Connection failed:",
				slog.Any("error", err))
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		logger.Info("Connected to:",
			slog.Any("peer", peer))
	}

	return node
}

func (node *nodeInputPort) Run() {
	node.logger.Info("Running p2p node")

}

func (node *nodeInputPort) Shutdown() {
	node.logger.Info("Gracefully shutting down p2p node")
	node.host.Close()
}
