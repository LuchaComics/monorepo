package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type CreateBlockDataService struct {
	config                 *config.Config
	logger                 *slog.Logger
	createBlockDataUseCase *usecase.CreateBlockDataUseCase
	getBlockDataUseCase    *usecase.GetBlockDataUseCase
}

func NewCreateBlockDataService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.CreateBlockDataUseCase,
	uc2 *usecase.GetBlockDataUseCase,
) *CreateBlockDataService {
	return &CreateBlockDataService{cfg, logger, uc1, uc2}
}

func (s *CreateBlockDataService) Execute(dataDir, hash, walletPassword string) (*domain.BlockData, error) {
	//
	// STEP 1: Valhashation.
	//

	e := make(map[string]string)
	if dataDir == "" {
		e["data_dir"] = "missing value"
	}
	if hash == "" {
		e["hash"] = "missing value"
	}
	if walletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating block create parameters",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	header := &domain.BlockHeader{}
	trans := make([]domain.BlockTransaction, 0)
	headerSignature := []byte{}
	validator := &domain.Validator{}

	//
	// STEP 3: Save to our database.
	//

	if err := s.createBlockDataUseCase.Execute(hash, header, headerSignature, trans, validator); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("hash", hash),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 4: Return the account.
	//

	return s.getBlockDataUseCase.Execute(hash)
}
