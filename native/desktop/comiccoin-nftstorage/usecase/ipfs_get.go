package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"

	domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type IPFSGetUseCase struct {
	logger   *slog.Logger
	ipfsRepo domain.IPFSRepository
}

func NewIPFSGetUseCase(logger *slog.Logger, r1 domain.IPFSRepository) *IPFSGetUseCase {
	return &IPFSGetUseCase{logger, r1}
}

func (uc *IPFSGetUseCase) Execute(ctx context.Context, cid string) ([]byte, error) {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)

	if cid == "" {
		e["cid"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed.",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Execute submitting to IPFS.
	//

	content, contentType, err := uc.ipfsRepo.Get(ctx, cid)
	if err != nil {
		uc.logger.Error("Failed adding to IPFS via file content",
			slog.Any("error", err))
		return nil, err
	}

	uc.logger.Debug("Fetched from IPFS successfully",
		slog.Any("cid", cid),
		slog.Any("contentType", contentType))

	return content, nil
}
