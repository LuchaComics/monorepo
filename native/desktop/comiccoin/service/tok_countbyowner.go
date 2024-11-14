package service

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type CountByOwnerTokenService struct {
	config                    *config.Config
	logger                    *slog.Logger
	countTokensByOwnerUseCase *usecase.CountTokensByOwnerUseCase
}

func NewCountByOwnerTokenService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.CountTokensByOwnerUseCase,
) *CountByOwnerTokenService {
	return &CountByOwnerTokenService{cfg, logger, uc1}
}

func (s *CountByOwnerTokenService) Execute(address *common.Address) (int, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating listing tokens by owner",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count the tokens by owner.
	//

	return s.countTokensByOwnerUseCase.Execute(address)
}
