package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type BroadcastSignedIssuedTokenDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedIssuedTokenDTORepository
}

func NewBroadcastSignedIssuedTokenDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedIssuedTokenDTORepository) *BroadcastSignedIssuedTokenDTOUseCase {
	return &BroadcastSignedIssuedTokenDTOUseCase{config, logger, repo}
}

func (uc *BroadcastSignedIssuedTokenDTOUseCase) Execute(ctx context.Context, ido *domain.SignedIssuedToken) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if ido == nil {
		e["token"] = "missing value"
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

	dto := &domain.SignedIssuedTokenDTO{
		IssuedToken: domain.IssuedToken{
			ID:          ido.ID,
			MetadataURI: ido.MetadataURI,
		},
		IssuedTokenSignatureBytes: ido.IssuedTokenSignatureBytes,
		Validator:                 ido.Validator,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.BroadcastToP2PNetwork(ctx, dto)
}
