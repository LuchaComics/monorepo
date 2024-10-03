package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type LastBlockDataHashDTOReceiveP2PRequestUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashDTORepository
}

func NewLastBlockDataHashDTOReceiveP2PRequestUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashDTORepository) *LastBlockDataHashDTOReceiveP2PRequestUseCase {
	return &LastBlockDataHashDTOReceiveP2PRequestUseCase{config, logger, repo}
}

func (uc *LastBlockDataHashDTOReceiveP2PRequestUseCase) Execute(ctx context.Context) (peer.ID, error) {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.ReceiveRequestFromNetwork(ctx)
}
