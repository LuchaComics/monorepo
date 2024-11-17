package service

import (
	"context"
	"log/slog"
	"math/big"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type TokenGetService struct {
	logger          *slog.Logger
	getTokenUseCase *usecase.GetTokenUseCase
}

func NewTokenGetService(
	logger *slog.Logger,
	uc1 *usecase.GetTokenUseCase,
) *TokenGetService {
	return &TokenGetService{logger, uc1}
}

func (s *TokenGetService) Execute(ctx context.Context, tokenID *big.Int) (*domain.Token, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if tokenID == nil {
		e["token_id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for getting token",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our token from our in-memory database if it exists.
	//

	token, err := s.getTokenUseCase.Execute(ctx, tokenID)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting token",
				slog.Any("token_id", tokenID),
				slog.Any("error", err))
			return nil, err
		}
	}
	if token != nil {
		return token, nil
	}

	return token, nil
}
