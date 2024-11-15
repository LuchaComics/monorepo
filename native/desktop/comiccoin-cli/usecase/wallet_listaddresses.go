package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListAllAddressesWalletUseCase struct {
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewListAllAddressesWalletUseCase(logger *slog.Logger, repo domain.WalletRepository) *ListAllAddressesWalletUseCase {
	return &ListAllAddressesWalletUseCase{logger, repo}
}

func (uc *ListAllAddressesWalletUseCase) Execute(ctx context.Context) ([]*common.Address, error) {
	addresses, err := uc.repo.ListAllAddresses(ctx)
	if err != nil {
		uc.logger.Error("Failed listing all by addresses",
			slog.Any("error", err))
		return nil, err
	}
	uc.logger.Debug("Addresses listed",
		slog.Any("addresses", addresses))
	return addresses, nil
}
