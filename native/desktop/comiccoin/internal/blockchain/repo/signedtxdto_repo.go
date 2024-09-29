package repo

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

const (
	rendezvousString = "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/signedtxdto"
)

type signedTransactionDTORepoImpl struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
}

func NewSignedTransactionDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.SignedTransactionDTORepository {
	impl := &signedTransactionDTORepoImpl{cfg, logger, libP2PNetwork}

	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), rendezvousString)

	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	go func() {
		impl.logger.Debug("waiting for peers...",
			slog.String("rendezvous_string", rendezvousString))

		for {
			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), rendezvousString, func(p peer.AddrInfo) error {
				impl.logger.Debug("connected",
					slog.Any("peer_id", p.ID),
					slog.Any("is_host", impl.libP2PNetwork.IsHostMode()),
					slog.String("rendezvous_string", rendezvousString))

				// TODO: Do something. (hint, setup pub-sub)

				// Return nil to indicate success
				return nil
			})
		}
	}()

	return impl
}

func (impl *signedTransactionDTORepoImpl) Broadcast(ctx context.Context, bd *domain.SignedTransactionDTO) error {

	return nil
}

func (impl *signedTransactionDTORepoImpl) Receive(ctx context.Context) (*domain.SignedTransactionDTO, error) {

	return nil, nil
}
