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

// Provider provides interface for abstracting P2P networking.
type LibP2PNetwork interface {
	// Returns the host node of the P2P network.
	GetHost() host.Host

	// Returns the pub-sub handler for this peer.
	GetPubSubSingletonInstance() *pubsub.PubSub

	// Returns whether the peer is in host mode or not.
	IsHostMode() bool

	// Advertise the peer's presence using a rendezvous string.
	AdvertiseWithRendezvousString(ctx context.Context, h host.Host, rendezvousString string)

	// Discover peers at a rendezvous string and connect to them.
	DiscoverPeersAtRendezvousString(ctx context.Context, h host.Host, rendezvousString string, peerConnectedFunc func(peer.AddrInfo) error)

	// Put data into the Kademlia DHT.
	PutDataToKademliaDHT(key string, bytes []byte) error

	// Get data from the Kademlia DHT.
	GetDataFromKademliaDHT(key string) ([]byte, error)

	// Remove data from the Kademlia DHT.
	RemoveDataFromKademliaDHT(key string) error

	// Close the P2P network connection.
	Close()
}

type peerProviderImpl struct {
	// The configuration for the P2P network.
	cfg *config.Config

	// The logger for the P2P network.
	logger *slog.Logger

	// The private key and public key for the peer's identity.
	identityPrivateKey crypto.PrivKey
	identityPublicKey  crypto.PubKey

	// A mutex to protect access to the subs and closed fields.
	mu sync.Mutex

	// The host node of the P2P network.
	host host.Host

	// Whether the peer is in host mode or not.
	isHostMode bool

	// The Kademlia DHT instance.
	kademliaDHT *dht.IpfsDHT

	// The routing discovery instance.
	routingDiscovery *routing.RoutingDiscovery

	// A map of connected peers, keyed by rendezvous string.
	peers map[string]map[peer.ID]*peer.AddrInfo

	// --- Publish-Subscriber variables below... ---

	// A flag to indicate whether the message queue broker has been closed.
	closed bool

	// A map of pub-sub topics, keyed by topic name.
	topics map[string]*pubsub.Topic

	// The gossip pub-sub instance.
	gossipPubSub *pubsub.PubSub
}

// NewLibP2PNetwork creates a new instance of the P2P network.
func NewLibP2PNetwork(cfg *config.Config, logger *slog.Logger, priv crypto.PrivKey, pub crypto.PubKey) LibP2PNetwork {
	// Create a new instance of the peer provider.
	impl := &peerProviderImpl{
		cfg:                cfg,
		logger:             logger,
		identityPrivateKey: priv,
		identityPublicKey:  pub,
		peers:              make(map[string]map[peer.ID]*peer.AddrInfo, 0),
	}

	// Create a new host node with a predictable identifier.
	h, err := impl.newHostWithPredictableIdentifier()
	if err != nil {
		log.Fatalf("failed to load host: %v", err)
	}
	impl.host = h

	// Create a new Kademlia DHT instance.
	kademliaDHT := impl.newKademliaDHT(context.Background())
	impl.kademliaDHT = kademliaDHT

	// Create a new routing discovery instance.
	routingDiscovery := drouting.NewRoutingDiscovery(impl.kademliaDHT)
	impl.routingDiscovery = routingDiscovery

	// Create a new gossip pub-sub instance.
	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		log.Fatalf("failed setting new gossip sub: %v", err)
	}
	impl.gossipPubSub = ps

	// Set up a notification handler for disconnected peers.
	impl.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			// Remove the disconnected peer from the list of connected peers.
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			for _, rendezvousPeers := range impl.peers {
				_, ok := rendezvousPeers[peerID]
				if ok {
					// Remove the peer from the host node's peerstore.
					h.Network().ClosePeer(peerID)
					h.Peerstore().RemovePeer(peerID)
					impl.kademliaDHT.RoutingTable().RemovePeer(peerID)

					// Remove the peer from the list of connected peers.
					delete(rendezvousPeers, peerID)
					impl.logger.Warn("deleted peer",
						slog.Any("rendezvousPeers", rendezvousPeers),
						slog.Any("rendezvousPeer", rendezvousPeers[peerID]))

					break
				}
			}
		},
	})

	return impl
}

// GetHost returns the host node of the P2P network.
func (p *peerProviderImpl) GetHost() host.Host {
	return p.host
}

// GetPubSubSingletonInstance returns the pub-sub handler for this peer.
func (p *peerProviderImpl) GetPubSubSingletonInstance() *pubsub.PubSub {
	return p.gossipPubSub
}

// IsHostMode returns whether the peer is in host mode or not.
func (p *peerProviderImpl) IsHostMode() bool {
	return p.isHostMode
}

// Close closes the P2P network connection.
func (impl *peerProviderImpl) Close() {
	// // Close the gossip pub-sub instance.
	// impl.gossipPubSub.Close()

	// Close the Kademlia DHT instance.
	impl.kademliaDHT.Close()

	// Close the host node.
	impl.host.Close()

	// Set the closed flag to true.
	impl.mu.Lock()
	impl.closed = true
	impl.mu.Unlock()

	// Log a message to indicate that the P2P network connection has been closed.
	impl.logger.Info("P2P network connection closed")
}
