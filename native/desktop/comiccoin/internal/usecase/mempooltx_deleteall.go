package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
)

type DeleteAllMempoolTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewDeleteAllMempoolTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.MempoolTransactionRepository) *DeleteAllMempoolTransactionUseCase {
	return &DeleteAllMempoolTransactionUseCase{config, logger, repo}
}

func (uc *DeleteAllMempoolTransactionUseCase) Execute() error {
	err := uc.repo.DeleteAll()
	if err != nil {
		uc.logger.Error("Failed deleting all mempool transactions",
			slog.Any("error", err))
		return err
	}
	return nil
}
