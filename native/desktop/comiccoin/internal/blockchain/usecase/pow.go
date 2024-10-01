package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type ProofOfWorkUseCase struct {
	config *config.Config
	logger *slog.Logger
}

func NewProofOfWorkUseCase(config *config.Config, logger *slog.Logger) *ProofOfWorkUseCase {
	return &ProofOfWorkUseCase{config, logger}
}

func (uc *ProofOfWorkUseCase) Execute(id string) (uint64, error) {
	//
	// STEP 1: Validation. (TODO: IMPL)
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	//TODO: IMPL

	return 0, nil
}
