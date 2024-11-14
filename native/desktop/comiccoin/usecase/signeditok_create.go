package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type CreateSignedIssuedTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedIssuedTokenRepository
}

func NewCreateSignedIssuedTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedIssuedTokenRepository) *CreateSignedIssuedTokenUseCase {
	return &CreateSignedIssuedTokenUseCase{config, logger, repo}
}

func (uc *CreateSignedIssuedTokenUseCase) Execute(sitok *domain.SignedIssuedToken) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if sitok.IssuedToken.MetadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if sitok.IssuedTokenSignatureBytes == nil {
		e["issued_token_signature_bytes"] = "missing value"
	}
	if sitok.Validator == nil {
		e["validator"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Upsert(sitok)
}
