package repo

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

type BlockDataDTORepo struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork

	rendezvousString string

	mu sync.Mutex

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo
}

func NewBlockDataDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.BlockDataDTORepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"

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
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			delete(impl.peers, peerID)

		},
	})

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
					slog.Any("peer_id", p.ID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (r *BlockDataDTORepo) UploadToNetwork(ctx context.Context, data *domain.BlockDataDTO) error {
	dataBytes, err := data.Serialize()
	if err != nil {
		return err
	}
	return r.libP2PNetwork.PutDataToKademliaDHT(data.Hash, dataBytes)

}

func (r *BlockDataDTORepo) DownloadFromNetwork(ctx context.Context, blockDataHash string) (*domain.BlockDataDTO, error) {
	dataBytes, err := r.libP2PNetwork.GetDataFromKademliaDHT(blockDataHash)
	if err != nil {
		return nil, err
	}
	data, err := domain.NewBlockDataDTOFromDeserialize(dataBytes)
	if err != nil {
		return nil, err
	}
	return data, nil
}
