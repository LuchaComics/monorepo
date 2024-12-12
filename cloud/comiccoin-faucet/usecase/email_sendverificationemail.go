package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/templatedemailer"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type SendUserVerificationEmailUseCase struct {
	config  *config.Configuration
	logger  *slog.Logger
	emailer templatedemailer.TemplatedEmailer
}

func NewSendUserVerificationEmailUseCase(config *config.Configuration, logger *slog.Logger, emailer templatedemailer.TemplatedEmailer) *SendUserVerificationEmailUseCase {
	return &SendUserVerificationEmailUseCase{config, logger, emailer}
}

func (uc *SendUserVerificationEmailUseCase) Execute(ctx context.Context, user *domain.User) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if user == nil {
		e["user"] = "User is missing value"
	} else {
		if user.FirstName == "" {
			e["first_name"] = "First name is required"
		}
		if user.Email == "" {
			e["email"] = "Email is required"
		}
		if user.EmailVerificationCode == "" {
			e["email_verification_code"] = "Email verification code is required"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Send email
	//

	return uc.emailer.SendUserVerificationEmail(ctx, user.Email, user.EmailVerificationCode, user.FirstName)
}
