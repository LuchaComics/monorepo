package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/usecase"
)

type GetTokenService struct {
	config          *config.Config
	logger          *slog.Logger
	getTokenUseCase *usecase.GetTokenUseCase
}

func NewGetTokenService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetTokenUseCase,
) *GetTokenService {
	return &GetTokenService{cfg, logger, uc1}
}

func (s *GetTokenService) Execute(tokenID uint64) (*domain.Token, error) {
	token, err := s.getTokenUseCase.Execute(tokenID)
	if err != nil {
		s.logger.Error("failed getting token",
			slog.Uint64("token_id", tokenID),
			slog.Any("error", err))
		return nil, err
	}
	return token, nil
}
