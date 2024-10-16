package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type BlockDataDTOSendP2PRequestUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewBlockDataDTOSendP2PRequestUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *BlockDataDTOSendP2PRequestUseCase {
	return &BlockDataDTOSendP2PRequestUseCase{config, logger, repo}
}

func (uc *BlockDataDTOSendP2PRequestUseCase) Execute(ctx context.Context, blockDataHash string) error {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.SendRequestToRandomPeer(ctx, blockDataHash)
}
