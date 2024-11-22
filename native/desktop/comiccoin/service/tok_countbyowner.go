package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type TokenCountByOwnerService struct {
	logger                    *slog.Logger
	countTokensByOwnerUseCase *usecase.CountTokensByOwnerUseCase
}

func NewTokenCountByOwnerService(
	logger *slog.Logger,
	uc1 *usecase.CountTokensByOwnerUseCase,
) *TokenCountByOwnerService {
	return &TokenCountByOwnerService{logger, uc1}
}

func (s *TokenCountByOwnerService) Execute(ctx context.Context, address *common.Address) (int64, error) {
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

	return s.countTokensByOwnerUseCase.Execute(ctx, address)
}
