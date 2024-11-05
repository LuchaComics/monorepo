package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ReceiveIssuedTokenDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IssuedTokenDTORepository
}

func NewReceiveIssuedTokenDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.IssuedTokenDTORepository) *ReceiveIssuedTokenDTOUseCase {
	return &ReceiveIssuedTokenDTOUseCase{config, logger, repo}
}

func (uc *ReceiveIssuedTokenDTOUseCase) Execute(ctx context.Context) (*domain.IssuedToken, []byte, *domain.Validator, error) {
	//
	// STEP 1:
	// Wait to receive from the P2P Network. It just takes one node to publish
	// the data and then we will receive it here.
	//

	dto, err := uc.repo.ReceiveFromP2PNetwork(ctx)
	if err != nil {
		uc.logger.Error("failed receiving issued token dto from network",
			slog.Any("error", err))
		return nil, nil, nil, err
	}
	if dto == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do is return nil.
		return nil, nil, nil, nil
	}
	if dto.Token == nil || dto.TokenSignatureBytes == nil || dto.Validator == nil {
		// Developer Note:
		// Any of these missing fields are grounds for immediate termination.
		return nil, nil, nil, nil
	}

	//
	// STEP 2:
	// Convert back to our signed trnsaction data-type and then perform simple
	// validation before returning it for this function.
	//

	ido := &domain.IssuedToken{
		ID:          dto.Token.ID,
		MetadataURI: dto.Token.MetadataURI,
	}

	e := make(map[string]string)
	if ido.MetadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received issued token",
			slog.Any("error", e))
		return nil, nil, nil, httperror.NewForBadRequest(&e)
	}

	uc.logger.Debug("Received issued token dto from network",
		slog.Any("token_id", ido.ID))

	return ido, dto.TokenSignatureBytes, dto.Validator, nil
}
