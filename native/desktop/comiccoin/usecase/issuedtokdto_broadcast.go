package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type BroadcastIssuedTokenDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IssuedTokenDTORepository
}

func NewBroadcastIssuedTokenDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.IssuedTokenDTORepository) *BroadcastIssuedTokenDTOUseCase {
	return &BroadcastIssuedTokenDTOUseCase{config, logger, repo}
}

func (uc *BroadcastIssuedTokenDTOUseCase) Execute(ctx context.Context, tok *domain.IssuedToken, sig []byte, val *domain.Validator) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if tok == nil {
		e["token"] = "missing value"
	}
	if sig == nil {
		e["token_signature_bytes"] = "missing value"
	}
	if val == nil {
		e["validator"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed.",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	// Get a signature of the token.

	dto := &domain.IssuedTokenDTO{
		Token:               tok,
		TokenSignatureBytes: sig,
		Validator:           val,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.BroadcastToP2PNetwork(ctx, dto)
}
