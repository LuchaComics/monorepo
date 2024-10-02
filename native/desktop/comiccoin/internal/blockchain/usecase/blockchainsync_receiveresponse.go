package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

type BlockchainSyncReceiveResponseUseCase struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	repo          domain.BlockDataRepository
}

func NewBlockchainSyncReceiveResponseUseCase(config *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork, repo domain.BlockDataRepository) *BlockchainSyncReceiveResponseUseCase {
	return &BlockchainSyncReceiveResponseUseCase{config, logger, libP2PNetwork, repo}
}

func (uc *BlockchainSyncReceiveResponseUseCase) Execute() error {
	return nil
}
