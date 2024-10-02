package repo

import (
	"context"
	"fmt"
	"log/slog"

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

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(blockDataDTOProtocolID, func(stream network.Stream) {
		// TODO: Handle the stream here
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

				// TODO: Figure out

				// host := libP2PNetwork.GetHost()
				// ctx := context.Background()
				//
				// stream, err := host.NewStream(ctx, p.ID, protocol.ID(blockDataDTOProtocolID))
				// if err != nil {
				// 	impl.logger.Error("Connection failed", slog.Any("error", err))
				// 	return err
				// }
				// _ = stream
				//
				// impl.logger.Debug("sync stream ready",
				// 	slog.Any("peer_id", p.ID),
				// 	slog.String("protocol_id", blockDataDTOProtocolID))

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

	// TODO: Figure out

	// host := r.libP2PNetwork.GetHost()
	// for _, peerInfo := range r.peers {
	// 	// peerInfo is a single record from the map
	// 	fmt.Println(peerInfo)
	//
	// 	stream, err := host.NewStream(ctx, peerInfo.ID, protocol.ID(blockDataDTOProtocolID))
	// 	if err != nil {
	// 		r.logger.Error("Connection failed", slog.Any("error", err))
	// 		return nil, err
	// 	}
	// 	fmt.Println("Todo: stream:", stream)
	//
	// 	break
	// }

	//TODO: IMPL.
	fmt.Println("Todo: ListLatestAfterHash")
	return nil, nil
}
