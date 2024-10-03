package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type LastBlockDataHashDTOReceiveP2PResponseUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashDTORepository
}

func NewLastBlockDataHashDTOReceiveP2PResponseUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashDTORepository) *LastBlockDataHashDTOReceiveP2PResponseUseCase {
	return &LastBlockDataHashDTOReceiveP2PResponseUseCase{config, logger, repo}
}

func (uc *LastBlockDataHashDTOReceiveP2PResponseUseCase) Execute(ctx context.Context) (domain.LastBlockDataHashDTO, error) {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.ReceiveResponseFromNetwork(ctx)
}
