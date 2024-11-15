package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type MempoolTransactionReceiveDTOFromNetworkService struct {
	config                          *config.Configuration
	logger                          *slog.Logger
	mempoolTransactionCreateUseCase *usecase.MempoolTransactionCreateUseCase
}

func NewMempoolTransactionReceiveDTOFromNetworkService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.MempoolTransactionCreateUseCase,
) *MempoolTransactionReceiveDTOFromNetworkService {
	return &MempoolTransactionReceiveDTOFromNetworkService{cfg, logger, uc}
}

func (s *MempoolTransactionReceiveDTOFromNetworkService) Execute(ctx context.Context, dto *domain.MempoolTransactionDTO) error {
	//
	// STEP 1: Validation.
	//

	if dto == nil {
		err := fmt.Errorf("Cannot have empty mempool transaction dto")
		s.logger.Warn("Validation failed for received",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2: Convert from data transfer object to internal data object.
	//

	mempoolTx := dto.ToIDO()

	//
	// STEP 3: Save to our database.
	//

	// DEVELOPERS NOTE:
	// What happens when we save the mempooltx submission? Our application
	// is constantly waiting to detect any new transactions in the database
	// and once we submit this current transaction, then the `Proof of
	// Authority` consensus mechanism will activate and execute. Therefore,
	// all we have to do from this service perspective is just save the
	// network submission.
	upsertErr := s.mempoolTransactionCreateUseCase.Execute(ctx, mempoolTx)
	if upsertErr != nil {
		s.logger.Warn("Validation failed for received",
			slog.Any("error", upsertErr))
		return upsertErr
	}
	return nil

}
