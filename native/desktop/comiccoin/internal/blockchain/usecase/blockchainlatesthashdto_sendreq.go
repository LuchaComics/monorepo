package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockchainLastestHashDTOSendP2PRequestUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashDTORepository
}

func NewBlockchainLastestHashDTOSendP2PRequestUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashDTORepository) *BlockchainLastestHashDTOSendP2PRequestUseCase {
	return &BlockchainLastestHashDTOSendP2PRequestUseCase{config, logger, repo}
}

func (uc *BlockchainLastestHashDTOSendP2PRequestUseCase) Execute(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.SendRequestToRandomPeer(ctx)
}
