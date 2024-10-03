package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type LastBlockDataHashDTOSendP2PRequestUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashDTORepository
}

func NewLastBlockDataHashDTOSendP2PRequestUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashDTORepository) *LastBlockDataHashDTOSendP2PRequestUseCase {
	return &LastBlockDataHashDTOSendP2PRequestUseCase{config, logger, repo}
}

func (uc *LastBlockDataHashDTOSendP2PRequestUseCase) Execute(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.SendRequestToRandomPeer(ctx)
}
