package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	pkg_domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type IPFSPinAddUseCase struct {
	logger   *slog.Logger
	ipfsRepo pkg_domain.IPFSRepository
}

func NewIPFSPinAddUseCase(logger *slog.Logger, r1 pkg_domain.IPFSRepository) *IPFSPinAddUseCase {
	return &IPFSPinAddUseCase{logger, r1}
}

func (uc *IPFSPinAddUseCase) Execute(fileContent []byte) (string, error) {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)

	if fileContent == nil {
		e["file_content"] = "missing value"
	}
	if len(e) != 0 {
		return "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Define our object.
	//

	return uc.ipfsRepo.AddViaFileContent(fileContent, true)
}
