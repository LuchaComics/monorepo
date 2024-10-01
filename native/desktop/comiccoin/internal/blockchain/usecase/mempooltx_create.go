package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateMempoolTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewCreateMempoolTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.MempoolTransactionRepository) *CreateMempoolTransactionUseCase {
	return &CreateMempoolTransactionUseCase{config, logger, repo}
}

func (uc *CreateMempoolTransactionUseCase) Execute(stx *domain.MempoolTransaction) error {
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
		e["value"] = "missing value"
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

	// //
	// // STEP 2: Create our strucutre.
	// //
	//
	// tx := domain.Transaction{
	// 	ChainID: chainID,
	// 	Nonce:   nonce,
	// 	From:    from,
	// 	To:      to,
	// 	Value:   value,
	// 	Tip:     tip,
	// 	Data:    data,
	// }
	// stx := &domain.MempoolTransaction{
	// 	Transaction: tx,
	// 	V:           v,
	// 	R:           r,
	// 	S:           s,
	// }

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(stx)
}
