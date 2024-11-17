package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type GetAccountService struct {
	logger            *slog.Logger
	getAccountUseCase *usecase.GetAccountUseCase
}

func NewGetAccountService(
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
) *GetAccountService {
	return &GetAccountService{logger, uc1}
}

func (s *GetAccountService) Execute(ctx context.Context, address *common.Address) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our account from our in-memory database if it exists.
	//

	account, err := s.getAccountUseCase.Execute(ctx, address)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting account",
				slog.Any("address", address),
				slog.Any("error", err))
			return nil, err
		}
	}
	if account != nil {
		return account, nil
	}

	return account, nil
}
