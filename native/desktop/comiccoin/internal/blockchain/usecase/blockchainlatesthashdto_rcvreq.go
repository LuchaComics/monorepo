package usecase

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockchainLastestHashDTOReceiveP2PRequestUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashDTORepository
}

func NewBlockchainLastestHashDTOReceiveP2PRequestUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashDTORepository) *BlockchainLastestHashDTOReceiveP2PRequestUseCase {
	return &BlockchainLastestHashDTOReceiveP2PRequestUseCase{config, logger, repo}
}

func (uc *BlockchainLastestHashDTOReceiveP2PRequestUseCase) Execute(ctx context.Context) (peer.ID, error) {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.ReceiveRequestFromNetwork(ctx)
}
