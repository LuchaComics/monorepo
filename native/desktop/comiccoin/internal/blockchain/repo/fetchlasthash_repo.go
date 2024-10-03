package repo

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type FetchLastHashRepo struct {
	config           *config.Config
	logger           *slog.Logger
	libP2PNetwork    p2p.LibP2PNetwork
	rendezvousString string

	doneCh chan bool

	node *Node

	mu sync.Mutex

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo
}

type FetchLastHashRepository interface {
	SendToRandomPeerInNetwork(hash string) (string, error)
}

func NewFetchLastHashRepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) FetchLastHashRepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"

	//
	// STEP 1
	// Initialize our instance
	//

	doneCh := make(chan bool, 1)

	impl := &FetchLastHashRepo{
		config:           cfg,
		logger:           logger,
		rendezvousString: rendezvousString,
		libP2PNetwork:    libP2PNetwork,
		doneCh:           doneCh,
	}

	//
	// STEP 2:
	// Create and advertise our `rendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	// This is like your friend telling you the location to meet you.
	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), impl.rendezvousString)

	//
	// STEP 3:
	// Load up all the stream handlers by this peer.
	//

	host := libP2PNetwork.GetHost()

	//
	// STEP 4:
	// In a peer-to-peer network we need to be away of when peers disconnect
	// our network, the following code will callback when a peer disconnects so
	// our repository can remove the peer from our records.
	//

	//Remove disconnected peer
	host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			delete(impl.peers, peerID)

		},
	})

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.

	peerNode := NewNode(host, doneCh)
	impl.node = peerNode

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...")

		for {

			//
			// STEP 6:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), impl.rendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 7
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				impl.logger.Debug("peer connected",
					slog.Any("peer_id", p.ID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (impl *FetchLastHashRepo) SendToRandomPeerInNetwork(hash string) (string, error) {
	yetAnotherRandomPeerID := impl.yetAnotherRandomPeerID()
	if yetAnotherRandomPeerID == "" {
		return "", fmt.Errorf("no peers connected")
	}

	requestID, err := impl.node.SendRequest(yetAnotherRandomPeerID, hash)
	if err != nil {
		return "", err
	}

	for {
		reponses := impl.node.FetchLastHashProtocol.GetResponse()
		response, ok := reponses[requestID]
		if !ok {
			time.Sleep(5 * time.Second)
			continue
		}
		if response != nil {
			return response.Hash, nil
		}
	}
}

func (impl *FetchLastHashRepo) ReceiveFromNetwork(hash string) (string, error) {
	yetAnotherRandomPeerID := impl.yetAnotherRandomPeerID()
	if yetAnotherRandomPeerID == "" {
		return "", fmt.Errorf("no peers connected")
	}

	requestID, err := impl.node.SendRequest(yetAnotherRandomPeerID, hash)
	if err != nil {
		return "", err
	}

	for {
		reponses := impl.node.FetchLastHashProtocol.GetResponse()
		response, ok := reponses[requestID]
		if !ok {
			time.Sleep(5 * time.Second)
			continue
		}
		if response != nil {
			return response.Hash, nil
		}
	}
}

func (r *FetchLastHashRepo) yetAnotherRandomPeerID() peer.ID {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get a list of peer IDs
	peerIDs := make([]peer.ID, 0, len(r.peers))
	for id := range r.peers {
		peerIDs = append(peerIDs, id)
	}

	// Select a random peer ID
	if len(peerIDs) == 0 {
		// Handle the case where there are no peers
		return ""
	}
	return peerIDs[rand.Intn(len(peerIDs))]
}
