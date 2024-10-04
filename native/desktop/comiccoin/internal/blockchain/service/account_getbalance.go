package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetAccountBalanceService struct {
	config                          *config.Config
	logger                          *slog.Logger
	getBlockchainLastestHashUseCase *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
}

func NewGetAccountBalanceService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainLastestHashUseCase,
	uc2 *usecase.GetBlockDataUseCase,
) *GetAccountBalanceService {
	return &GetAccountBalanceService{cfg, logger, uc1, uc2}
}

func (s *GetAccountBalanceService) Execute(account *domain.Account) (uint64, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if account == nil {
		e["account"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed getting account balance",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//TODO: IMPL.

	return 0, nil
}
