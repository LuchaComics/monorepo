package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockchainLastestHashDTOReceiveP2PResponseUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashDTORepository
}

func NewBlockchainLastestHashDTOReceiveP2PResponseUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashDTORepository) *BlockchainLastestHashDTOReceiveP2PResponseUseCase {
	return &BlockchainLastestHashDTOReceiveP2PResponseUseCase{config, logger, repo}
}

func (uc *BlockchainLastestHashDTOReceiveP2PResponseUseCase) Execute(ctx context.Context) (domain.BlockchainLastestHashDTO, error) {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.ReceiveResponseFromNetwork(ctx)
}
