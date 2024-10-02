package usecase

import (
	"log/slog"

	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p/p2pmessagedto"
)

type BlockchainSyncSendRequestUseCase struct {
	config             *config.Config
	logger             *slog.Logger
	libP2PNetwork      p2p.LibP2PNetwork
	p2pMessengeDTORepo p2pmessagedto.P2PMessageDTORepository
}

func NewBlockchainSyncSendRequestUseCase(config *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) *BlockchainSyncSendRequestUseCase {
	impl := &BlockchainSyncSendRequestUseCase{
		config:        config,
		logger:        logger,
		libP2PNetwork: libP2PNetwork}

	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"
	protocolID := protocol.ID("/sync/1.0.0")

	p2pMessengeDTORepo := p2pmessagedto.NewP2PMessageDTORepo(logger, libP2PNetwork, rendezvousString, protocolID)
	impl.p2pMessengeDTORepo = p2pMessengeDTORepo

	return impl
}

func (uc *BlockchainSyncSendRequestUseCase) Execute() error {
	return nil
}
