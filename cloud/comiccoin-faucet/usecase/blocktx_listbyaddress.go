package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/ethereum/go-ethereum/common"
)

//
// Copied from `github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase`
//

type ListBlockTransactionsByAddressUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockTransactionsByAddressUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockTransactionsByAddressUseCase {
	return &ListBlockTransactionsByAddressUseCase{config, logger, repo}
}

func (uc *ListBlockTransactionsByAddressUseCase) Execute(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	data, err := uc.repo.ListBlockTransactionsByAddress(ctx, address)
	if err != nil {
		uc.logger.Error("failed listing block transactions by address",
			slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
