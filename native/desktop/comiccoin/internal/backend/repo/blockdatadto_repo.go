package repo

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p/protocol/simple"
)

type BlockDataDTORepo struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	dtoProtocol   simple.SimpleDTOProtocol

	rendezvousString string

	mu sync.Mutex

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo
}

func NewBlockDataDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.BlockDataDTORepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain/blockdatadto"

	//
	// STEP 1
	// Initialize our instance
	//

	impl := &BlockDataDTORepo{
		config:           cfg,
		logger:           logger,
		libP2PNetwork:    libP2PNetwork,
		rendezvousString: rendezvousString,
		peers:            make(map[peer.ID]*peer.AddrInfo, 0),
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

	h := libP2PNetwork.GetHost()

	//
	// STEP 4:
	// In a peer-to-peer network we need to be away of when peers disconnect
	// our network, the following code will callback when a peer disconnects so
	// our repository can remove the peer from our records.
	//

	//Remove disconnected peer
	h.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected",
				slog.Any("remote_peer_id", peerID),
				slog.Any("local_peer_id", h.ID()),
				slog.String("dto", "blockdatadto"),
			)
			delete(impl.peers, peerID)

		},
	})

	dtoProtocol := simple.NewSimpleDTOProtocol(logger, h, "/blockdatadto/req/1.0.0", "/blockdatadto/resp/1.0.0")
	impl.dtoProtocol = dtoProtocol

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

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), h, impl.rendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 7
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				impl.logger.Debug("peer connected",
					slog.String("dto", "blockdatadto"),
					slog.Any("remote_peer_id", p.ID),
					slog.Any("local_peer_id", h.ID()))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
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

func (impl *BlockDataDTORepo) SendRequestToRandomPeer(ctx context.Context, blockDataHash string) error {
	randomPeerID := impl.randomPeerID()
	if randomPeerID == "" {
		return fmt.Errorf("no peers connected")
	}
	if err := impl.dtoProtocol.SendRequest(randomPeerID, []byte(blockDataHash)); err != nil {
		return err
	}

	return nil
}

func (impl *BlockDataDTORepo) ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, string, error) {
	req, err := impl.dtoProtocol.WaitAndReceiveRequest(ctx)
	if err != nil {
		impl.logger.Error("failed receiving download request from network",
			slog.Any("error", err))
		return "", "", err
	}
	blockDataHash := string(req.Content)
	return req.FromPeerID, blockDataHash, nil
}

func (impl *BlockDataDTORepo) SendResponseToPeer(ctx context.Context, peerID peer.ID, data *domain.BlockDataDTO) error {
	dataBytes, err := cbor.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize stream message dto: %v", err)
	}

	if err := impl.dtoProtocol.SendResponse(peerID, dataBytes); err != nil {
		impl.logger.Error("failed sending upload response from network",
			slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *BlockDataDTORepo) ReceiveResponseFromNetwork(ctx context.Context) (*domain.BlockDataDTO, error) {
	resp, err := impl.dtoProtocol.WaitAndReceiveResponse(ctx)
	if err != nil {
		return nil, err
	}

	// Variable we will use to return.
	dto := &domain.BlockDataDTO{}

	if err := cbor.Unmarshal(resp.Content, &dto); err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}

	return dto, nil
}
