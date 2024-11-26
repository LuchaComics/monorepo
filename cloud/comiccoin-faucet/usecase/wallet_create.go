package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type CreateWalletUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewCreateWalletUseCase(config *config.Configuration, logger *slog.Logger, repo domain.WalletRepository) *CreateWalletUseCase {
	return &CreateWalletUseCase{config, logger, repo}
}

func (uc *CreateWalletUseCase) Execute(ctx context.Context, address *common.Address, keystoreBytes []byte, label string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if keystoreBytes == nil {
		e["keystore_bytes"] = "missing value"
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
		Address:       address,
		KeystoreBytes: keystoreBytes,
		Label:         label,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, wallet)
}
