package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"

	domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type IPFSPinAddUseCase struct {
	logger   *slog.Logger
	ipfsRepo domain.IPFSRepository
}

func NewIPFSPinAddUseCase(logger *slog.Logger, r1 domain.IPFSRepository) *IPFSPinAddUseCase {
	return &IPFSPinAddUseCase{logger, r1}
}

func (uc *IPFSPinAddUseCase) Execute(data []byte) (string, error) {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)

	if data == nil {
		e["data"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed.",
			slog.Any("e", e))
		return "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Execute submitting to IPFS.
	//

	cid, err := uc.ipfsRepo.AddViaFileContent(data, true)
	if err != nil {
		uc.logger.Error("Failed adding to IPFS via file content",
			slog.Any("error", err))
		return "", err
	}

	uc.logger.Debug("Submitted to IPFS successfully",
		slog.Any("cid", cid))

	return cid, nil
}
