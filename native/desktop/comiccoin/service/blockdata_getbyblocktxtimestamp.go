package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type GetByBlockTransactionTimestampService struct {
	logger                                *slog.Logger
	GetByBlockTransactionTimestampUseCase *usecase.GetByBlockTransactionTimestampUseCase
}

func NewGetByBlockTransactionTimestampService(
	logger *slog.Logger,
	uc1 *usecase.GetByBlockTransactionTimestampUseCase,
) *GetByBlockTransactionTimestampService {
	return &GetByBlockTransactionTimestampService{logger, uc1}
}

func (s *GetByBlockTransactionTimestampService) Execute(ctx context.Context, timestamp uint64) (*domain.BlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if timestamp == 0 {
		e["timestamp"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get
	//

	blockData, err := s.GetByBlockTransactionTimestampUseCase.Execute(ctx, timestamp)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting block data",
				slog.Any("timestamp", timestamp),
				slog.Any("error", err))
			return nil, err
		}
	}
	if blockData != nil {
		return blockData, nil
	}

	return nil, nil
}
