package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/common"
)

// DEPRECATED
type UpsertTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewUpsertTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *UpsertTokenUseCase {
	return &UpsertTokenUseCase{config, logger, repo}
}

func (uc *UpsertTokenUseCase) Execute(
	id uint64,
	owner *common.Address,
	metadataURI string,
	// nonce uint64,
) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if owner == nil {
		e["owner"] = "missing value"
	}
	if metadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upsert our structure.
	//

	token := &domain.Token{
		ID:          id,
		Owner:       owner,
		MetadataURI: metadataURI,
		// Nonce:       nonce,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(token)
}
