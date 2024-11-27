package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type UserGetByVerificationCodeUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserRepository
}

func NewUserGetByVerificationCodeUseCase(config *config.Configuration, logger *slog.Logger, repo domain.UserRepository) *UserGetByVerificationCodeUseCase {
	return &UserGetByVerificationCodeUseCase{config, logger, repo}
}

func (uc *UserGetByVerificationCodeUseCase) Execute(ctx context.Context, verificationCode string) (*domain.User, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if verificationCode == "" {
		e["verification_code"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for get by verification",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Get from database.
	//

	return uc.repo.GetByVerificationCode(ctx, verificationCode)
}
