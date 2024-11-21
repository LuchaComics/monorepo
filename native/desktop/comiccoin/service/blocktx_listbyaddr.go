package service

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ListBlockTransactionsByAddressService struct {
	logger                                *slog.Logger
	listBlockTransactionsByAddressUseCase *usecase.ListBlockTransactionsByAddressUseCase
}

func NewListBlockTransactionsByAddressService(
	logger *slog.Logger,
	uc1 *usecase.ListBlockTransactionsByAddressUseCase,
) *ListBlockTransactionsByAddressService {
	return &ListBlockTransactionsByAddressService{logger, uc1}
}

func (s *ListBlockTransactionsByAddressService) Execute(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: List the tokens by owner.
	//

	return s.listBlockTransactionsByAddressUseCase.Execute(ctx, address)
}
