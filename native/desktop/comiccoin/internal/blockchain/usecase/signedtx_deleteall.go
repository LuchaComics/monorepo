package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type DeleteAllSignedTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedTransactionRepository
}

func NewDeleteAllSignedTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedTransactionRepository) *DeleteAllSignedTransactionUseCase {
	return &DeleteAllSignedTransactionUseCase{config, logger, repo}
}

func (uc *DeleteAllSignedTransactionUseCase) Execute() error {
	err := uc.repo.DeleteAll()
	if err != nil {
		uc.logger.Error("Failed deleting all signed transactions",
			slog.Any("error", err))
		return err
	}
	return nil
}
