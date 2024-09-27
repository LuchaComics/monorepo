package peer

import (
	"context"
	"log"
	"log/slog"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/config"
)

// Provider provides interface for abstracting P2P netowrking.
type LibP2PNetwork interface {
	// Returns your peer's host.
	GetHost() host.Host

	// Your peer advertises a rendezvous point and waits for other peers to join.
	// When other peers make a successful connection to your peer, then method
	// will send the connected peer through the channel for you to grab.
	DiscoverPeersChannel(ctx context.Context, h host.Host, rendezvousString string) <-chan peer.AddrInfo

	Close()
}

type peerProviderImpl struct {
	cfg    *config.Config
	logger *slog.Logger

	identityPrivateKey crypto.PrivKey
	identityPublicKey  crypto.PubKey

	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery

	// mu field is used to protect access to the subs and closed fields using a mutex.
	mu sync.Mutex
}

// NewLibP2PNetwork constructor that returns the default P2P connected instance.
func NewLibP2PNetwork(cfg *config.Config, logger *slog.Logger, priv crypto.PrivKey, pub crypto.PubKey) LibP2PNetwork {
	ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	impl := &peerProviderImpl{
		cfg:                cfg,
		logger:             logger,
		identityPrivateKey: priv,
		identityPublicKey:  pub,
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

func (impl *peerProviderImpl) Close() {

}
