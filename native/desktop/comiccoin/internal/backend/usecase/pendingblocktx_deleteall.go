package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
)

type DeleteAllPendingBlockTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.PendingBlockTransactionRepository
}

func NewDeleteAllPendingBlockTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.PendingBlockTransactionRepository) *DeleteAllPendingBlockTransactionUseCase {
	return &DeleteAllPendingBlockTransactionUseCase{config, logger, repo}
}

func (uc *DeleteAllPendingBlockTransactionUseCase) Execute() error {
	err := uc.repo.DeleteAll()
	if err != nil {
		uc.logger.Error("Failed deleting all pending block transactions",
			slog.Any("error", err))
		return err
	}
	return nil
}
