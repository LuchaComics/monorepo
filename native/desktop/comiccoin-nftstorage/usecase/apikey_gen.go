package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/password"
)

type GenerateAPIKeyUseCase struct {
	logger   *slog.Logger
	password password.Provider
	jwt      jwt.Provider
}

func NewGenerateAPIKeyUseCase(logger *slog.Logger, passp password.Provider, jwtp jwt.Provider) *GenerateAPIKeyUseCase {
	return &GenerateAPIKeyUseCase{logger, passp, jwtp}
}

type NFTStoreAppCredentials struct {
	SecretString string
	APIKey       string
}

func (uc *GenerateAPIKeyUseCase) Execute(chainID uint16) (*NFTStoreAppCredentials, error) {

	// Generate hash for the secret.
	randomSecretStr, err := uc.password.GenerateSecureRandomString(64)
	if err != nil {
		uc.logger.Error("Failed generating secure random string error",
			slog.Any("error", err))
		return nil, err
	}
	randomSecretHash, err := uc.password.GenerateHashFromPassword(randomSecretStr)
	if err != nil {
		uc.logger.Error("Failed hashing error",
			slog.Any("error", err))
		return nil, err
	}

	// Generate our one-time API key and attach it to the response. What is
	// important here is that we share the plaintext secret to the user to
	// keep but we do not keep the plaintext value in our system, we only
	// keep the hash, so we keep the value safe.
	apiKeyPayload := fmt.Sprintf("%v@%v", chainID, randomSecretStr)
	atExpiry := 250 * 24 * time.Hour // Duration: 250 years.
	apiKey, _, err := uc.jwt.GenerateJWTToken(apiKeyPayload, atExpiry)
	if err != nil {
		uc.logger.Error("jwt generate pairs error",
			slog.Any("err", err))
		return nil, err
	}

	return &NFTStoreAppCredentials{
		SecretString: randomSecretHash,
		APIKey:       apiKey,
	}, nil

}
