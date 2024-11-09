package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type GetWalletUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewGetWalletUseCase(config *config.Configuration, logger *slog.Logger, repo domain.WalletRepository) *GetWalletUseCase {
	return &GetWalletUseCase{config, logger, repo}
}

func (uc *GetWalletUseCase) Execute(ctx context.Context, address *common.Address) (*domain.Wallet, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting wallet",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByAddress(ctx, address)
}
