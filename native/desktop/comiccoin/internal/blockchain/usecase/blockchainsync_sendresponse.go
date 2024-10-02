package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

type BlockchainSyncSendResponseUseCase struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	repo          domain.BlockDataRepository
}

func NewBlockchainSyncSendResponseUseCase(config *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork, repo domain.BlockDataRepository) *BlockchainSyncSendResponseUseCase {
	return &BlockchainSyncSendResponseUseCase{config, logger, libP2PNetwork, repo}
}

func (uc *BlockchainSyncSendResponseUseCase) Execute() error {
	return nil
}
