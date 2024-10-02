package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	blockDataDTORendezvousString = "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"
	blockDataDTOProtocolID       = "/sync/1.0.0"
)

type BlockDataDTORepo struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	// The list of connected peers.
	peers map[peer.ID]*peer.AddrInfo
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

	// Remove disconnected peer
	host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			delete(impl.peers, peerID)
		},
	})

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(blockDataDTOProtocolID, func(stream network.Stream) {
		// Handle incoming streams
		switch stream.Protocol() {
		case blockDataDTOProtocolID:
			// Handle ListLatestAfterHash request
			impl.handleListLatestAfterHashRequest(stream)
		default:
			// Handle unknown protocol
			log.Println("Unknown protocol:", stream.Protocol())
		}
	})

	//
	// STEP 4:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...",
			slog.String("protocol_id", blockDataDTOProtocolID))

		for {

			//
			// STEP 5:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), blockDataDTORendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 6
				// Connect our peer with this peer and keep a record of it.
				//

				impl.logger.Debug("setting up blockdata dto sync stream...",
					slog.Any("peer_id", p.ID),
					slog.String("protocol_id", blockDataDTOProtocolID))

				impl.peers[p.ID] = &p

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (r *BlockDataDTORepo) ListLatestAfterHash(ctx context.Context, afterBlockDataHash string) ([]*domain.BlockDataDTO, error) {
	if len(r.peers) == 0 {
		r.logger.Warn("No peers")
		return []*domain.BlockDataDTO{}, nil
	}

	xxx, err := r.ListLatestAfterHashV2(ctx, "")
	if err != nil {
		return nil, err
	}

	//TODO: IMPL.
	fmt.Println("Todo: ListLatestAfterHash", xxx)
	return nil, nil
}

type ListLatestAfterHashRequest struct {
	AfterBlockDataHash string `json:"after_block_data_hash"`
}

type ListLatestAfterHashResponse struct {
	BlockDataDTOs []*domain.BlockDataDTO `json:"block_data_dtos"`
}

func (r *BlockDataDTORepo) handleListLatestAfterHashRequest(stream network.Stream) {
	// Read the request from the stream
	var req ListLatestAfterHashRequest
	err := json.NewDecoder(stream).Decode(&req)
	if err != nil {
		log.Println("Error decoding request:", err)
		return
	}

	// Process the request
	blockDataDTOs, err := r.ListLatestAfterHashV2(context.Background(), req.AfterBlockDataHash)
	if err != nil {
		log.Println("Error processing request:", err)
		return
	}

	// Send the response back to the client
	resp := ListLatestAfterHashResponse{BlockDataDTOs: blockDataDTOs}
	err = json.NewEncoder(stream).Encode(resp)
	if err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}

func (r *BlockDataDTORepo) ListLatestAfterHashV2(ctx context.Context, afterBlockDataHash string) ([]*domain.BlockDataDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to a random peer
	peer := r.peers[r.randomPeerID()]
	if peer == nil {
		return nil, errors.New("no peers connected")
	}

	// Create a new stream to the peer
	stream, err := r.libP2PNetwork.GetHost().NewStream(ctx, peer.ID, blockDataDTOProtocolID)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	// Write the request to the stream
	req := ListLatestAfterHashRequest{AfterBlockDataHash: afterBlockDataHash}
	err = json.NewEncoder(stream).Encode(req)
	if err != nil {
		stream.Reset() // Reset the stream in case of an error
		return nil, err
	}

	// Read the response from the stream
	var resp ListLatestAfterHashResponse
	err = json.NewDecoder(stream).Decode(&resp)
	if err != nil {
		stream.Reset() // Reset the stream in case of an error
		return nil, err
	}

	return resp.BlockDataDTOs, nil
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
