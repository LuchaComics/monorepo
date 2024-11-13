package tokens

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/big"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func GetTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get account details",
		Run: func(cmd *cobra.Command, args []string) {
			doRunGetToken()
		},
	}

	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The token ID value to lookup the account by")
	cmd.MarkFlagRequired("token-id")

	return cmd
}

func doRunGetToken() {
	// Common
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)

	// Repository
	tokRepo := repo.NewTokenRepo(cfg, logger, dbClient)

	// // Use-case
	getTokenUseCase := usecase.NewGetTokenUseCase(
		cfg,
		logger,
		tokRepo,
	)

	// // Service
	tokenRetrieveService := service.NewTokenRetrieveService(
		logger,
		getTokenUseCase,
	)

	////
	//// Start the transaction.
	////
	ctx := context.Background()

	session, err := dbClient.StartSession()
	if err != nil {
		logger.Error("start session error",
			slog.Any("error", err))
		log.Fatalf("Failed executing: %v\n", err)
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		tokenID, ok := new(big.Int).SetString(flagTokenID, 10)
		if !ok {
			return nil, fmt.Errorf("Failed converting to big.Int: %v", flagTokenID)
		}
		logger.Debug("Getting token...",
			slog.Any("token_id_str", flagTokenID),
			slog.Any("token_id", tokenID.Uint64()))

		tok, err := tokenRetrieveService.Execute(sessCtx, tokenID)
		if err != nil {
			return nil, err
		}
		if tok == nil {
			err := errors.New("Token does not exist")
			logger.Error("Failed getting token",
				slog.Any("tokenID", tokenID.Uint64()),
				slog.Any("error", err))
			return nil, err
		}

		return tok, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed creating account: %v\n", err)
	}

	tok := res.(*domain.Token)

	logger.Debug("Token retrieved",
		slog.Any("id", tok.GetID()),
		slog.Any("nonce", tok.GetNonce()),
		slog.String("metadata_uri", tok.MetadataURI),
		slog.String("owner", tok.Owner.Hex()),
	)
}
