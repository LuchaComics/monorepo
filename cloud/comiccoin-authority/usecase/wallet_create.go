package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type CreateWalletUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewCreateWalletUseCase(config *config.Configuration, logger *slog.Logger, repo domain.WalletRepository) *CreateWalletUseCase {
	return &CreateWalletUseCase{config, logger, repo}
}

func (uc *CreateWalletUseCase) Execute(ctx context.Context, address *common.Address, filepath string, label string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if filepath == "" {
		e["filepath"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating new wallet",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	wallet := &domain.Wallet{
		Address:  address,
		FilePath: filepath,
		Label:    label,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, wallet)
}
