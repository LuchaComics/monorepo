package service

import (
	"context"
	"log/slog"
	"math/big"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type TokenRetrieveService struct {
	logger          *slog.Logger
	getTokenUseCase *usecase.GetTokenUseCase
}

func NewTokenRetrieveService(
	logger *slog.Logger,
	uc1 *usecase.GetTokenUseCase,
) *TokenRetrieveService {
	return &TokenRetrieveService{logger, uc1}
}

func (s *TokenRetrieveService) Execute(ctx context.Context, id *big.Int) (*domain.Token, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == nil {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for getting token",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our token from our in-memory database if it exists.
	//

	token, err := s.getTokenUseCase.Execute(ctx, id)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting token",
				slog.Any("id", id),
				slog.Any("error", err))
			return nil, err
		}
	}
	if token != nil {
		return token, nil
	}

	return token, nil
}
