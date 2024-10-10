package usecase

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreatePendingBlockTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.PendingBlockTransactionRepository
}

func NewCreatePendingBlockTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.PendingBlockTransactionRepository) *CreatePendingBlockTransactionUseCase {
	return &CreatePendingBlockTransactionUseCase{config, logger, repo}
}

func (uc *CreatePendingBlockTransactionUseCase) Execute(stx *domain.PendingBlockTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if stx.ChainID != uc.config.Blockchain.ChainID {
		e["chain_id"] = "wrong blockchain used"
	}
	// Nonce - skip validating.
	if stx.From == nil {
		e["from"] = "missing value"
	}
	if stx.To == nil {
		e["to"] = "missing value"
	}
	if stx.Value <= 0 {
		// DEVELOPERS NOTE:
		// Only `coin` type transactions need their value verified while the
		// `nft` type transactions can have zero value.
		if stx.Type == domain.TransactionTypeCoin {
			e["value"] = "missing value"
		}
	}
	if stx.Type == "" {
		e["type"] = "missing value"
	} else {
		var validType bool = false
		if stx.Type == domain.TransactionTypeCoin {
			validType = true
		}
		if stx.Type == domain.TransactionTypeNFT {
			validType = true

			if stx.TokenMetadataURI == "" {
				e["token_metadata_uri"] = "missing value"
			}
		}
		if validType == false {
			e["type"] = fmt.Sprintf("incorrect value: %v", stx.Type)
		}
	}
	// Tip - skip validating.
	// Data - skip validating.
	if stx.V == nil {
		e["v"] = "missing value"
	}
	if stx.R == nil {
		e["r"] = "missing value"
	}
	if stx.S == nil {
		e["s"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Upsert(stx)
}
