package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type TokenListByOwnerService struct {
	logger                   *slog.Logger
	listTokensByOwnerUseCase *usecase.ListTokensByOwnerUseCase
}

func NewTokenListByOwnerService(
	logger *slog.Logger,
	uc1 *usecase.ListTokensByOwnerUseCase,
) *TokenListByOwnerService {
	return &TokenListByOwnerService{logger, uc1}
}

func (s *TokenListByOwnerService) Execute(ctx context.Context, address *common.Address) ([]*domain.Token, error) {
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

	return s.listTokensByOwnerUseCase.Execute(ctx, address)
}
