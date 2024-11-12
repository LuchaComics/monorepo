package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type GetWalletUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewGetWalletUseCase(config *config.Config, logger *slog.Logger, repo domain.WalletRepository) *GetWalletUseCase {
	return &GetWalletUseCase{config, logger, repo}
}

func (uc *GetWalletUseCase) Execute(address *common.Address) (*domain.Wallet, error) {
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

	return uc.repo.GetByAddress(address)
}
