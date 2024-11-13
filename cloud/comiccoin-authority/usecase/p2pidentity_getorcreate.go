package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetOrCreateLibP2PNetworkPeerUniqueIdentifierUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.LibP2PNetworkPeerUniqueIdentifierRepository
}

func NewGetOrCreateLibP2PNetworkPeerUniqueIdentifierUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.LibP2PNetworkPeerUniqueIdentifierRepository,
) *GetOrCreateLibP2PNetworkPeerUniqueIdentifierUseCase {
	return &GetOrCreateLibP2PNetworkPeerUniqueIdentifierUseCase{config, logger, repo}
}

func (uc *GetOrCreateLibP2PNetworkPeerUniqueIdentifierUseCase) Execute(ctx context.Context, label string) (*domain.LibP2PNetworkPeerUniqueIdentifier, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if label == "" {
		e["label"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get or insert.
	//

	return uc.repo.GetOrCreate(ctx, label)
}
