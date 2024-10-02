package repo

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"math/rand"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	blockDataDTORendezvousString = "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"
	blockDataDTOProtocolID       = "/sync/1.0.0"
)

type BlockDataDTORepo struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo

	// The list of connected peers with a direct stream with tem.
	streams map[peer.ID]network.Stream
}

func NewBlockDataDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) *BlockDataDTORepo {
	//
	// STEP 1
	// Initialize our instance
	//

	impl := &BlockDataDTORepo{
		config:        cfg,
		logger:        logger,
		libP2PNetwork: libP2PNetwork,
		peers:         make(map[peer.ID]*peer.AddrInfo, 0),
		streams:       make(map[peer.ID]network.Stream, 0),
	}

	//
	// STEP 2:
	// Create and advertise our `blockDataDTORendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	// This is like your friend telling you the location to meet you.
	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), blockDataDTORendezvousString)

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

			impl.logger.Warn("stream closed", slog.Any("peer_id", peerID))
			stream, ok := impl.streams[peerID]
			if ok {
				stream.Close()
				delete(impl.streams, peerID)

			}
		},
	})

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(blockDataDTOProtocolID, func(stream network.Stream) {
		// Handle incoming streams
		switch stream.Protocol() {
		case blockDataDTOProtocolID:
			impl.streams[host.ID()] = stream
		default:
			// Handle unknown protocol
			log.Fatalf("Unknown protocol: %v", stream.Protocol())
		}
	})

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...",
			slog.String("protocol_id", blockDataDTOProtocolID))

		for {

			//
			// STEP 6:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), blockDataDTORendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 7
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				ctx := context.Background()
				stream, err := host.NewStream(ctx, p.ID, protocol.ID(blockDataDTOProtocolID))
				if err != nil {
					// logger.Warn("Connection failed", slog.Any("error", err))
					return err
				} else {
					impl.streams[p.ID] = stream
				}

				impl.logger.Debug("peer connected",
					slog.Any("peer_id", p.ID),
					slog.String("protocol_id", blockDataDTOProtocolID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (r *BlockDataDTORepo) ListLatestAfterHash(ctx context.Context, afterBlockDataHash string) ([]*domain.BlockDataDTO, error) {
	randomPeerID := r.randomPeerID()
	if randomPeerID == "" {
		return nil, nil
	}

	return nil, nil
}

func (r *BlockDataDTORepo) randomPeerID() peer.ID {
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

func (r *BlockDataDTORepo) getRandomStream() (network.Stream, error) {
	peerID := r.randomPeerID()
	if peerID == "" {
		return nil, errors.New("no valid peers available")
	}

	s, _ := r.streams[peerID]
	return s, nil
}

func (r *BlockDataDTORepo) getRandomPeer() (*peer.AddrInfo, error) {
	peerID := r.randomPeerID()
	if peerID == "" {
		return nil, errors.New("no valid peers available")
	}

	// Connect to a random peer
	peer, _ := r.peers[peerID]
	if peer == nil {
		return nil, errors.New("no peers connected")
	}
	return peer, nil
}

type ListLatestAfterHashRequest struct {
	AfterBlockDataHash string `json:"after_block_data_hash"`
}

type ListLatestAfterHashResponse struct {
	BlockDataDTOs []*domain.BlockDataDTO `json:"block_data_dtos"`
}
