package p2p

import (
	"context"
	"log"
	"log/slog"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
)

// Provider provides interface for abstracting P2P netowrking.
type LibP2PNetwork interface {
	// Returns your peer's host.
	GetHost() host.Host

	// Return the pub-sub handler for this peer.
	GetPubSubSingletonInstance() *pubsub.PubSub

	IsHostMode() bool

	AdvertiseWithRendezvousString(ctx context.Context, h host.Host, rendezvousString string)

	DiscoverPeersAtRendezvousString(ctx context.Context, h host.Host, rendezvousString string, peerConnectedFunc func(peer.AddrInfo) error)

	// Your peer advertises a rendezvous point and waits for other peers to join.
	// When other peers make a successful connection to your peer, then method
	// will send the connected peer through the channel for you to grab.

	Close()
}

type peerProviderImpl struct {
	cfg    *config.Config
	logger *slog.Logger

	identityPrivateKey crypto.PrivKey
	identityPublicKey  crypto.PubKey

	// mu field is used to protect access to the subs and closed fields using a mutex.
	mu sync.Mutex

	host             host.Host
	isHostMode       bool
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery

	// The list of connected peers.
	peers map[string]map[peer.ID]*peer.AddrInfo

	// --- Publish-Subscriber variables below... ---

	// The quit field is a channel that is closed when the `message queue broker` is closed, allowing goroutines that are blocked on the channel to unblock and exit.
	quit chan struct{}

	// The closed field is a flag that indicates whether the `message queue broker` has been closed.
	closed bool

	topics       map[string]*pubsub.Topic
	gossipPubSub *pubsub.PubSub
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
		peers:              make(map[string]map[peer.ID]*peer.AddrInfo, 0),
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

	// Load up the gossip pub-sub
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Fatalf("failed setting new gossip sub: %v", err)
	}
	impl.gossipPubSub = ps

	//Remove disconnected peer
	impl.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			for _, rendezvousPeers := range impl.peers {
				_, ok := rendezvousPeers[peerID]
				if ok {
					// STEP 1:
					// Fetch our record.
					peer := rendezvousPeers[peerID]

					// STEP 2:
					// Remove our peer from our libp2p networking
					h.Network().ClosePeer(peer.ID)
					h.Peerstore().RemovePeer(peer.ID)
					impl.kademliaDHT.RoutingTable().RemovePeer(peer.ID)

					// STEP 2:
					// Remove the peer from our library
					delete(rendezvousPeers, peerID)
					impl.logger.Warn("deleted peer",
						slog.Any("rendezvousPeers", rendezvousPeers),
						slog.Any("rendezvousPeer", rendezvousPeers[peerID]))

					break
				}
			}
			//
		},
	})

	return impl
}

func (p *peerProviderImpl) GetHost() host.Host {
	return p.host
}

func (p *peerProviderImpl) GetPubSubSingletonInstance() *pubsub.PubSub {
	return p.gossipPubSub
}

func (p *peerProviderImpl) IsHostMode() bool {
	return p.isHostMode
}

func (impl *peerProviderImpl) Close() {

}
