package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetIdentityKeyService struct {
	config                *config.Config
	logger                *slog.Logger
	getIdentityKeyUseCase *usecase.GetIdentityKeyUseCase
}

func NewGetIdentityKeyService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.GetIdentityKeyUseCase,
) *GetIdentityKeyService {
	return &GetIdentityKeyService{cfg, logger, uc}
}

func (s *GetIdentityKeyService) Execute(id string) (*domain.IdentityKey, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed getting identity because validation failed",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	return s.getIdentityKeyUseCase.Execute(id)
}
