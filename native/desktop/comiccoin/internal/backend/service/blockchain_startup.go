package service

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
)

type BlockchainStartupService struct {
	config                            *config.Config
	logger                            *slog.Logger
	initAccountsFromBlockchainService *InitAccountsFromBlockchainService
	initBlockDataService              *InitBlockDataService
}

func NewBlockchainStartupService(
	cfg *config.Config,
	logger *slog.Logger,
	s1 *InitAccountsFromBlockchainService,
	s2 *InitBlockDataService,
) *BlockchainStartupService {
	return &BlockchainStartupService{cfg, logger, s1, s2}
}

func (s *BlockchainStartupService) Execute() error {

	// Load up the accounts into the in-memory storage before loading
	// the application because our accounts are only stored in memory
	// and not on disk.
	s.logger.Info("Loading accounts into memory...")
	if err := s.initAccountsFromBlockchainService.Execute(); err != nil {
		log.Fatalf("failed executing accounts initialization: %v\n", err)
	}
	s.logger.Info("Accounts database ready.")

	s.logger.Info("Initializing blockchain database...")
	if err := s.initBlockDataService.Execute(); err != nil {
		log.Fatalf("failed executing blockdata initialization: %v\n", err)
	}
	s.logger.Info("Blockchain database ready")

	return nil
}
