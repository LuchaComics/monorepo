package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/ethereum/go-ethereum/common"
)

// Struct represents the use case of looking up the previous token record and
// only update the record if the new nonce value is greater then or equal to
// the previous old nonce value.
//
// Why this use case (UC)? This UC is useful when we traverse the blockchain
// from most recent to the genesis because this UC will only save/update the
// database with the most recent account transaction (since most recent
// transactions have higher nonce values) and therefore ignore the previous
// transactions. We do this because the `token` database only shows the most
// recent tokens and their current owners, not the history of ownership.
type UpsertTokenIfPreviousTokenNonceGTEUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewUpsertTokenIfPreviousTokenNonceGTEUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *UpsertTokenIfPreviousTokenNonceGTEUseCase {
	return &UpsertTokenIfPreviousTokenNonceGTEUseCase{config, logger, repo}
}

func (uc *UpsertTokenIfPreviousTokenNonceGTEUseCase) Execute(
	id uint64,
	owner *common.Address,
	metadataURI string,
	nonce uint64,
) error {
	//
	// STEP 1:
	// Validation.
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
	// STEP 2:
	// Lookup previous record.
	//

	previousToken, err := uc.repo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed getting token by id", slog.Any("error", err))
		return err
	}

	//
	// CASE 1 OF 2:
	// Previous token D.N.E. therefore all we have to do is insert token.
	//
	if previousToken == nil {
		token := &domain.Token{
			ID:          id,
			Owner:       owner,
			MetadataURI: metadataURI,
			Nonce:       nonce,
		}
		return uc.repo.Upsert(token)
	}

	//
	// CASE 2 OF 2:
	// Previous record exists, so we must preform our logic.
	//

	//
	// STEP 3:
	// Compare `nonce` values and if nonce is not GTE then exit this function.
	//

	isGTE := nonce >= previousToken.Nonce

	if !isGTE {
		return nil
	}

	//
	// STEP 4:
	// Else nonce is GTE so we will upset our token record in the database.
	//

	token := &domain.Token{
		ID:          id,
		Owner:       owner,
		MetadataURI: metadataURI,
		Nonce:       nonce,
	}
	return uc.repo.Upsert(token)
}