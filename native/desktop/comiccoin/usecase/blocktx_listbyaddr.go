package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type ListBlockTransactionsByAddressUseCase struct {
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockTransactionsByAddressUseCase(logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockTransactionsByAddressUseCase {
	return &ListBlockTransactionsByAddressUseCase{logger, repo}
}

func (uc *ListBlockTransactionsByAddressUseCase) Execute(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.ListBlockTransactionsByAddress(ctx, address)
}