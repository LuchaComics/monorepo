package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type ReceiveSignedTransactionDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedTransactionDTORepository
}

func NewReceiveSignedTransactionDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedTransactionDTORepository) *ReceiveSignedTransactionDTOUseCase {
	return &ReceiveSignedTransactionDTOUseCase{config, logger, repo}
}

func (uc *ReceiveSignedTransactionDTOUseCase) Execute(ctx context.Context) (*domain.SignedTransaction, error) {
	//
	// STEP 1:
	// Wait to receive from the P2P Network. It just takes one node to publish
	// the data and then we will receive it here.
	//

	dto, err := uc.repo.ReceiveFromP2PNetwork(ctx)
	if err != nil {
		uc.logger.Warn("failed receiving signed transaction dto from network",
			slog.Any("error", err))
		return nil, err
	}
	if dto == nil {
		uc.logger.Warn("failed receiving signed transaction dto from network",
			slog.Any("error", "dto does not exist"))
		return nil, fmt.Errorf("received dto does not exist")
	}

	//
	// STEP 2:
	// Convert back to our signed trnsaction data-type and then perform simple
	// validation before returning it for this function.
	//

	ido := &domain.SignedTransaction{
		Transaction: dto.Transaction,
		V:           dto.V,
		R:           dto.R,
		S:           dto.S,
	}

	e := make(map[string]string)
	if ido.ChainID != uc.config.Blockchain.ChainID {
		e["chain_id"] = "wrong blockchain used"
	}
	// Nonce - skip validating.
	if ido.From == nil {
		e["from"] = "missing value"
	}
	if ido.To == nil {
		e["to"] = "missing value"
	}
	if ido.Value <= 0 {
		e["value"] = "missing value"
	}
	// Tip - skip validating.
	// Data - skip validating.
	if ido.S == nil {
		e["s"] = "missing value"
	}
	if ido.R == nil {
		e["r"] = "missing value"
	}
	if ido.V == nil {
		e["v"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	return ido, nil
}
