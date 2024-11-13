package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type MempoolTransactionCreateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewMempoolTransactionCreateUseCase(config *config.Configuration, logger *slog.Logger, repo domain.MempoolTransactionRepository) *MempoolTransactionCreateUseCase {
	return &MempoolTransactionCreateUseCase{config, logger, repo}
}

func (uc *MempoolTransactionCreateUseCase) Execute(ctx context.Context, mempoolTx *domain.MempoolTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if mempoolTx.ChainID != uc.config.Blockchain.ChainID {
		e["chain_id"] = "wrong blockchain used"
	}
	// Nonce - skip validating.
	if mempoolTx.From == nil {
		e["from"] = "missing value"
	}
	if mempoolTx.To == nil {
		e["to"] = "missing value"
	}
	if mempoolTx.Value <= 0 {
		// DEVELOPERS NOTE:
		// Only `coin` type transactions need their value verified while the
		// `token` type transactions can have zero value.
		if mempoolTx.Type == domain.TransactionTypeCoin {
			e["value"] = "missing value"
		}
	}
	if mempoolTx.Type == "" {
		e["type"] = "missing value"
	} else {
		var validType bool = false
		if mempoolTx.Type == domain.TransactionTypeCoin {
			validType = true
		}
		if mempoolTx.Type == domain.TransactionTypeToken {
			validType = true

			if mempoolTx.TokenMetadataURI == "" {
				e["token_metadata_uri"] = "missing value"
			}
		}
		if validType == false {
			e["type"] = fmt.Sprintf("incorrect value: %v", mempoolTx.Type)
		}
	}
	// Tip - skip validating.
	// Data - skip validating.
	if mempoolTx.V == nil {
		e["v"] = "missing value"
	}
	if mempoolTx.R == nil {
		e["r"] = "missing value"
	}
	if mempoolTx.S == nil {
		e["s"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Validate
	//

	if err := mempoolTx.Validate(uc.config.Blockchain.ChainID, true); err != nil {
		uc.logger.Warn("Validation failed for create",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, mempoolTx)
}
