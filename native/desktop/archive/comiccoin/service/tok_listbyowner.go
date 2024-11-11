package service

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ListByOwnerTokenService struct {
	config                   *config.Config
	logger                   *slog.Logger
	listTokensByOwnerUseCase *usecase.ListTokensByOwnerUseCase
}

func NewListByOwnerTokenService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ListTokensByOwnerUseCase,
) *ListByOwnerTokenService {
	return &ListByOwnerTokenService{cfg, logger, uc1}
}

func (s *ListByOwnerTokenService) Execute(address *common.Address, limit int) ([]*domain.Token, error) {
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
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: List the tokens by owner.
	//

	return s.listTokensByOwnerUseCase.Execute(address)
}
