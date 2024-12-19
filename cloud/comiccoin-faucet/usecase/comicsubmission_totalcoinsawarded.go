package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type ComicSubmissionTotalCoinsAwardedUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.ComicSubmissionRepository
}

func NewComicSubmissionTotalCoinsAwardedUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.ComicSubmissionRepository,
) *ComicSubmissionTotalCoinsAwardedUseCase {
	return &ComicSubmissionTotalCoinsAwardedUseCase{config, logger, repo}
}

func (uc *ComicSubmissionTotalCoinsAwardedUseCase) Execute(ctx context.Context) (uint64, error) {
	return uc.repo.TotalCoinsAwarded(ctx)
}
