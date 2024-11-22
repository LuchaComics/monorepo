package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type ListWithLimitBlockTransactionsByAddressUseCase struct {
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListWithLimitBlockTransactionsByAddressUseCase(logger *slog.Logger, repo domain.BlockDataRepository) *ListWithLimitBlockTransactionsByAddressUseCase {
	return &ListWithLimitBlockTransactionsByAddressUseCase{logger, repo}
}

func (uc *ListWithLimitBlockTransactionsByAddressUseCase) Execute(ctx context.Context, address *common.Address, limit int64) ([]*domain.BlockTransaction, error) {
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
		uc.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.ListWithLimitForBlockTransactionsByAddress(ctx, address, limit)
}

//
