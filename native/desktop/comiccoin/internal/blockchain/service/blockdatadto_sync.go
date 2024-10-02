package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type SyncBlockDataDTOService struct {
	config                                 *config.Config
	logger                                 *slog.Logger
	listLatestBlockDataAfterHashDTOUseCase *usecase.ListLatestBlockDataAfterHashDTOUseCase
}

func NewSyncBlockDataDTOService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.ListLatestBlockDataAfterHashDTOUseCase,
) *SyncBlockDataDTOService {
	return &SyncBlockDataDTOService{cfg, logger, uc}
}

func (s *SyncBlockDataDTOService) Execute(ctx context.Context) error {
	// hash := "" // Blank means start from genesis
	// data, err := s.listLatestBlockDataAfterHashDTOUseCase.Execute(ctx, hash)
	// if err != nil {
	// 	return err
	// }
	// if data == nil {
	// 	return nil
	// }
	// s.logger.Debug("----===>", slog.Any("data", data))
	return nil
}
