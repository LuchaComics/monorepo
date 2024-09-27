package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetBlockDataService struct {
	config              *config.Config
	logger              *slog.Logger
	getBlockDataUseCase *usecase.GetBlockDataUseCase
}

func NewGetBlockDataService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.GetBlockDataUseCase,
) *GetBlockDataService {
	return &GetBlockDataService{cfg, logger, uc}
}

func (s *GetBlockDataService) Execute(hash string) (*domain.BlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if hash == "" {
		e["hash"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	return s.getBlockDataUseCase.Execute(hash)
}
