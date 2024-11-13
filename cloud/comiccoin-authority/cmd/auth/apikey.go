package auth

import (
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/password"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func GenerateAPIKeyCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genapikey",
		Short: "Commands used to create a new API key for this service",
		Run: func(cmd *cobra.Command, args []string) {
			doGenerateAPIKeyCmd()
		},
	}

	return cmd
}

func doGenerateAPIKeyCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	// Developers Note:
	// To create a `` then run the following in your console:
	// `openssl rand -hex 64`

	// Misc
	cfg := config.NewProvider()
	logger := logger.NewProvider()
	passp := password.NewProvider()
	jwtp := jwt.NewProvider(cfg)
	apiKeyGenUseCase := usecase.NewGenerateAPIKeyUseCase(logger, passp, jwtp)

	//
	// STEP 2
	// Generate our applications credentials.
	//

	creds, err := apiKeyGenUseCase.Execute(cfg.Blockchain.ChainID)
	if err != nil {
		log.Fatalf("Failed to generate API key: %v\n", err)
	}

	//
	// STEP 3
	// Print to console.
	//

	logger.Info("Credentials created",
		slog.Any("api_key", creds.APIKey),
		slog.Any("auth_secret", creds.SecretString),
	)

}
