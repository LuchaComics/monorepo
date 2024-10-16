package usecase

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type CreateWalletUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewCreateWalletUseCase(config *config.Config, logger *slog.Logger, repo domain.WalletRepository) *CreateWalletUseCase {
	return &CreateWalletUseCase{config, logger, repo}
}

func (uc *CreateWalletUseCase) Execute(address *common.Address, filepath string) error {
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
		Filepath: filepath,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(wallet)
}
