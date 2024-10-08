package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type InitBlockDataService struct {
	config                          *config.Config
	logger                          *slog.Logger
	loadGenesisBlockDataUseCase     *usecase.LoadGenesisBlockDataUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
	createBlockDataUseCase          *usecase.CreateBlockDataUseCase
	setBlockchainLastestHashUseCase *usecase.SetBlockchainLastestHashUseCase
}

func NewInitBlockDataService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.LoadGenesisBlockDataUseCase,
	uc2 *usecase.GetBlockDataUseCase,
	uc3 *usecase.CreateBlockDataUseCase,
	uc4 *usecase.SetBlockchainLastestHashUseCase,
) *InitBlockDataService {
	return &InitBlockDataService{cfg, logger, uc1, uc2, uc3, uc4}
}

func (s *InitBlockDataService) Execute() error {
	//
	// STEP 1
	// Check to see if we have our genesis block in our database, and if so
	// then our blockchain is ready
	//

	genesis, err := s.getBlockDataUseCase.Execute(signature.ZeroHash)
	if err != nil {
		s.logger.Error("Failed getting genesis block", slog.Any("error", err))
		return err
	}
	if genesis != nil {
		// Blockchain initialization completed, exit this function.
		return nil
	}

	//
	// STEP 2:
	// Create our genesis block.
	//

	gbd, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting genesis block", slog.Any("error", err))
		return err
	}
	if err := s.createBlockDataUseCase.Execute(gbd.Hash, gbd.Header, gbd.Trans); err != nil {
		s.logger.Error("Failed creating genesis block", slog.Any("error", err))
		return err
	}

	//
	// STEP 3:
	// Set our latest hash to equal the genesis block data.
	//

	if err := s.setBlockchainLastestHashUseCase.Execute(gbd.Hash); err != nil {
		s.logger.Error("Failed creating genesis block", slog.Any("error", err))
		return err
	}

	return nil
}
