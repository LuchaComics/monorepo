package service

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ListWithLimitBlockTransactionsByAddressService struct {
	logger                                         *slog.Logger
	listWithLimitBlockTransactionsByAddressUseCase *usecase.ListWithLimitBlockTransactionsByAddressUseCase
}

func NewListWithLimitBlockTransactionsByAddressService(
	logger *slog.Logger,
	uc1 *usecase.ListWithLimitBlockTransactionsByAddressUseCase,
) *ListWithLimitBlockTransactionsByAddressService {
	return &ListWithLimitBlockTransactionsByAddressService{logger, uc1}
}

func (s *ListWithLimitBlockTransactionsByAddressService) Execute(ctx context.Context, address *common.Address, limit int64) ([]*domain.BlockTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if limit == 0 {
		e["limit"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: List the tokens by owner.
	//

	return s.listWithLimitBlockTransactionsByAddressUseCase.Execute(ctx, address, limit)
}
