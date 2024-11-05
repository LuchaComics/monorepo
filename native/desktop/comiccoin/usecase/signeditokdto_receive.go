package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ReceiveSignedIssuedTokenDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedIssuedTokenDTORepository
}

func NewReceiveSignedIssuedTokenDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedIssuedTokenDTORepository) *ReceiveSignedIssuedTokenDTOUseCase {
	return &ReceiveSignedIssuedTokenDTOUseCase{config, logger, repo}
}

func (uc *ReceiveSignedIssuedTokenDTOUseCase) Execute(ctx context.Context) (*domain.SignedIssuedToken, error) {
	//
	// STEP 1:
	// Wait to receive from the P2P Network. It just takes one node to publish
	// the data and then we will receive it here.
	//

	dto, err := uc.repo.ReceiveFromP2PNetwork(ctx)
	if err != nil {
		uc.logger.Error("failed receiving issued token dto from network",
			slog.Any("error", err))
		return nil, err
	}
	if dto == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do is return nil.
		return nil, nil
	}

	//
	// STEP 2:
	// Convert back to our signed trnsaction data-type and then perform simple
	// validation before returning it for this function.
	//

	ido := &domain.SignedIssuedToken{
		IssuedToken: domain.IssuedToken{
			ID:          dto.ID,
			MetadataURI: dto.MetadataURI,
		},
		V: dto.V,
		R: dto.R,
		S: dto.S,
	}

	e := make(map[string]string)
	if ido.MetadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received issued token",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	uc.logger.Debug("Received issued token dto from network",
		slog.Any("token_id", ido.ID))

	return ido, nil
}
