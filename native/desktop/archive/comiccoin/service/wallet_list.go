package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type WalletListService struct {
	config               *config.Config
	logger               *slog.Logger
	listAllWalletUseCase *usecase.ListAllWalletUseCase
}

func NewWalletListService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.ListAllWalletUseCase,
) *WalletListService {
	return &WalletListService{cfg, logger, uc}
}

func (s *WalletListService) Execute() ([]*domain.Wallet, error) {
	return s.listAllWalletUseCase.Execute()
}
