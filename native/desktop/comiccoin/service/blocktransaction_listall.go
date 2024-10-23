package service

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ListAllBlockTransactionService struct {
	config                                  *config.Config
	logger                                  *slog.Logger
	listAllBlockTransactionByAddressUseCase *usecase.ListAllBlockTransactionByAddressUseCase
}

func NewListAllBlockTransactionService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ListAllBlockTransactionByAddressUseCase,
) *ListAllBlockTransactionService {
	return &ListAllBlockTransactionService{cfg, logger, uc1}
}

func (s *ListAllBlockTransactionService) Execute(address *common.Address, limit int) ([]*domain.BlockTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating list All block transaction",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the data.
	//

	return s.listAllBlockTransactionByAddressUseCase.Execute(address)
}
