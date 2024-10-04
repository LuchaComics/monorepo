package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetAccountService struct {
	config            *config.Config
	logger            *slog.Logger
	getAccountUseCase *usecase.GetAccountUseCase
}

func NewGetAccountService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.GetAccountUseCase,
) *GetAccountService {
	return &GetAccountService{cfg, logger, uc}
}

func (s *GetAccountService) Execute(id string) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	return s.getAccountUseCase.Execute(id)
}
