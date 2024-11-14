package cmd

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/password"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/usecase"
)

// Command line argument flags
var (
	flatHMACSecret string
)

func GenerateAPIKeyCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genapikey",
		Short: "Commands used to create a new API key for this service",
		Run: func(cmd *cobra.Command, args []string) {
			doGenerateAPIKeyCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flatHMACSecret, "hmac-secret", "", "The HMAC secret to apply in our app")
	cmd.MarkFlagRequired("hmac-secret")

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
	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID: constants.ComicCoinChainID,
		},
		App: config.AppConfig{
			DirPath:     flagDataDir,
			HMACSecret:  []byte(flatHMACSecret),
			HTTPAddress: flagListenHTTPAddress,
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
