package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/ethereum/go-ethereum/common"
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

func (s *TokenListByOwnerService) Execute(ctx context.Context, ownerAddr *common.Address) ([]*domain.Token, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if ownerAddr == nil {
		e["ownerAddr"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed list tokens by owner",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our token from our in-memory database if it exists.
	//

	token, err := s.listTokensByOwnerUseCase.Execute(ctx, ownerAddr)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting token",
				slog.Any("owner", ownerAddr),
				slog.Any("error", err))
			return nil, err
		}
	}
	if token != nil {
		return token, nil
	}

	return token, nil
}
