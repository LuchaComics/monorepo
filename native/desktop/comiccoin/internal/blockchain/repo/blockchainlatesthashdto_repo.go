package repo

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	protoc "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo/protocol/blockchainlatesthashdto"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

type BlockchainLastestHashDTORepo struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	dtoProtocol   protoc.BlockchainLatestHashDTOProtocol

	rendezvousString string

	mu sync.Mutex

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo
}

func NewBlockchainLastestHashDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.BlockchainLastestHashDTORepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/lastblockdatahashdto"

	//
	// STEP 1
	// Initialize our instance
	//

	impl := &BlockchainLastestHashDTORepo{
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
				slog.Any("peer_id", peerID),
				slog.String("dto", "lastblockdatahash"),
			)
			delete(impl.peers, peerID)

		},
	})

	//
	dtoProtocol := protoc.NewBlockchainLatestHashDTOProtocol(logger, h)
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
					slog.String("dto", "lastblockdatahash"),
					slog.Any("peer_id", p.ID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (impl *BlockchainLastestHashDTORepo) SendRequestToRandomPeer(ctx context.Context) error {
	randomPeerID := impl.randomPeerID()
	if randomPeerID == "" {
		return fmt.Errorf("no peers connected")
	}

	impl.logger.Debug("SendRequestToRandomPeer: running...")

	// Note: Send empty request because we don't want anything.
	_, err := impl.dtoProtocol.SendRequest(randomPeerID, []byte(""))
	if err != nil {
		return err
	}

	impl.logger.Debug("SendRequestToRandomPeer: done")

	return nil
}

func (impl *BlockchainLastestHashDTORepo) ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error) {

	impl.logger.Debug("ReceiveRequestFromNetwork: running...")
	req, err := impl.dtoProtocol.WaitAndReceiveRequest(ctx)
	if err != nil {
		impl.logger.Error("failed receiving request from network", slog.Any("error", err))
		return "", err
	}
	impl.logger.Debug("ReceiveRequestFromNetwork: done")
	impl.logger.Debug("ReceiveRequestFromNetwork: results", slog.Any("req", req))
	return req.FromPeerID, nil
}

func (impl *BlockchainLastestHashDTORepo) SendResponseToPeer(ctx context.Context, peerID peer.ID, data domain.BlockchainLastestHashDTO) error {
	dataBytes := []byte(data)
	_, err := impl.dtoProtocol.SendResponse(peerID, dataBytes)
	if err != nil {
		return err
	}
	return nil

}

func (impl *BlockchainLastestHashDTORepo) ReceiveResponseFromNetwork(ctx context.Context) (domain.BlockchainLastestHashDTO, error) {
	resp, err := impl.dtoProtocol.WaitAndReceiveResponse(ctx)
	if err != nil {
		return "", err
	}

	hash := domain.BlockchainLastestHashDTO(resp.Content)
	return hash, nil
}

func (r *BlockchainLastestHashDTORepo) randomPeerID() peer.ID {
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
