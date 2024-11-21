package cmd

import (
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/password"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/usecase"
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

	dataDir := config.GetEnvString("COMICCOIN_NFTSTORAGE_APP_DATA_DIRECTORY", true)
	hmacSecretKey := config.GetSecureBytesEnv("COMICCOIN_NFTSTORAGE_APP_HMAC_SECRET", true)

	// Developers Note:
	// To create a `` then run the following in your console:
	// `openssl rand -hex 64`

	// Misc
	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID: constants.ComicCoinChainID,
		},
		App: config.AppConfig{
			DirPath:    dataDir,
			HMACSecret: hmacSecretKey,
		},
	}
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
		slog.Any("secret", creds.SecretString),
	)

}
