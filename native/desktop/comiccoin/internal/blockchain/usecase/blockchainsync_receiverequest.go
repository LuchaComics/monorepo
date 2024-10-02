package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

type BlockchainSyncReceiveRequestUseCase struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
}

func NewBlockchainSyncReceiveRequestUseCase(config *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) *BlockchainSyncReceiveRequestUseCase {
	return &BlockchainSyncReceiveRequestUseCase{config, logger, libP2PNetwork}
}

func (uc *BlockchainSyncReceiveRequestUseCase) Execute() error {
	return nil
}
