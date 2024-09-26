package peer

import (
	"context"
	"log"
	"log/slog"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// Provider provides interface for abstracting P2P netowrking.
type Provider interface {
	// Returns your peer's host.
	GetHost() host.Host

	// Your peer advertises a rendezvous point and waits for other peers to join.
	// When other peers make a successful connection to your peer, then method
	// will send the connected peer through the channel for you to grab.
	DiscoverPeersChannel(ctx context.Context, h host.Host, rendezvousString string) <-chan peer.AddrInfo
}

type peerProviderImpl struct {
	cfg           *config.Config
	logger        *slog.Logger
	keypairStorer keypair_ds.KeypairStorer

	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery
}

// NewProvider constructor that returns the default P2P connected instance.
func NewProvider(cfg *config.Config, logger *slog.Logger, kp keypair_ds.KeypairStorer) Provider {
	ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	impl := &peerProviderImpl{
		cfg:           cfg,
		logger:        logger,
		keypairStorer: kp,
	}

	// Run the code which will setup our peer-to-peer node in either `host mode`
	// or `dial mode`.
	h, err := impl.newHostWithPredictableIdentifier()
	if err != nil {
		log.Fatalf("failed to load host: %v", err)
	}
	impl.host = h

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	//
	// Source: https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/chat.go#L112
	kademliaDHT := impl.newKademliaDHT(ctx)
	impl.kademliaDHT = kademliaDHT

	routingDiscovery := drouting.NewRoutingDiscovery(impl.kademliaDHT)
	impl.routingDiscovery = routingDiscovery

	// peerCh := impl.DiscoverPeersChannel(ctx, h, cfg.Peer.RendezvousString)
	// peer := <-peerCh

	return impl
}

func (p *peerProviderImpl) GetHost() host.Host {
	return p.host
}
