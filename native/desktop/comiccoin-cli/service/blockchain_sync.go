package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

type BlockchainSyncWithBlockchainAuthorityService struct {
	logger                                               *slog.Logger
	getGenesisBlockDataUseCase                           *usecase.GetGenesisBlockDataUseCase
	upsertGenesisBlockDataUseCase                        *usecase.UpsertGenesisBlockDataUseCase
	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase *auth_usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase
	getBlockchainStateUseCase                            *usecase.GetBlockchainStateUseCase
	upsertBlockchainStateUseCase                         *usecase.UpsertBlockchainStateUseCase
	getBlockchainStateDTOFromBlockchainAuthorityUseCase  *auth_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase
	getBlockDataUseCase                                  *usecase.GetBlockDataUseCase
	upsertBlockDataUseCase                               *usecase.UpsertBlockDataUseCase
	getBlockDataDTOFromBlockchainAuthorityUseCase        *auth_usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase
}

func NewBlockchainSyncWithBlockchainAuthorityService(
	logger *slog.Logger,
	uc1 *usecase.GetGenesisBlockDataUseCase,
	uc2 *usecase.UpsertGenesisBlockDataUseCase,
	uc3 *auth_usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase,
	uc4 *usecase.GetBlockchainStateUseCase,
	uc5 *usecase.UpsertBlockchainStateUseCase,
	uc6 *auth_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase,
	uc7 *usecase.GetBlockDataUseCase,
	uc8 *usecase.UpsertBlockDataUseCase,
	uc9 *auth_usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase,
) *BlockchainSyncWithBlockchainAuthorityService {
	return &BlockchainSyncWithBlockchainAuthorityService{logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9}
}

func (s *BlockchainSyncWithBlockchainAuthorityService) Execute(ctx context.Context, chainID uint16) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chainID"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	return nil
}
