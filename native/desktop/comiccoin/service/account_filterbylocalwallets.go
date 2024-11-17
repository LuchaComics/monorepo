package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type AccountListingByLocalWalletsService struct {
	logger                           *slog.Logger
	listAllAddressesWalletUseCase    *usecase.ListAllAddressesWalletUseCase
	accountsFilterByAddressesUseCase *usecase.AccountsFilterByAddressesUseCase
}

func NewAccountListingByLocalWalletsService(
	logger *slog.Logger,
	uc1 *usecase.ListAllAddressesWalletUseCase,
	uc2 *usecase.AccountsFilterByAddressesUseCase,
) *AccountListingByLocalWalletsService {
	return &AccountListingByLocalWalletsService{logger, uc1, uc2}
}

func (s *AccountListingByLocalWalletsService) Execute(ctx context.Context) ([]*domain.Account, error) {
	myAccountAddresses, err := s.listAllAddressesWalletUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("Failed listing all wallet addresses",
			slog.Any("error", err))
		return nil, err
	}

	// If we have no local wallets then return an empty list.
	if myAccountAddresses == nil || len(myAccountAddresses) == 0 {
		return make([]*domain.Account, 0), nil
	}

	myAccounts, err := s.accountsFilterByAddressesUseCase.Execute(ctx, myAccountAddresses)
	if err != nil {
		s.logger.Error("Failed listing all my accounts",
			slog.Any("error", err))
		return nil, err
	}

	return myAccounts, nil
}
