package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateIdentityKeyService struct {
	config                   *config.Config
	logger                   *slog.Logger
	createIdentityKeyUseCase *usecase.CreateIdentityKeyUseCase
	getIdentityKeyUseCase    *usecase.GetIdentityKeyUseCase
}

func NewCreateIdentityKeyService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.CreateIdentityKeyUseCase,
	uc2 *usecase.GetIdentityKeyUseCase,
) *CreateIdentityKeyService {
	return &CreateIdentityKeyService{cfg, logger, uc1, uc2}
}

func (s *CreateIdentityKeyService) Execute(id string) (*domain.IdentityKey, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new identity key",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Save to our database.
	//

	if err := s.createIdentityKeyUseCase.Execute(id); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 3: Return the account.
	//

	return s.getIdentityKeyUseCase.Execute(id)
}
