package genesis

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func NewGenesistCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Initializes a new blockchain by creating a genesis block",
		Run: func(cmd *cobra.Command, args []string) {
			doRunNewAccount()
		},
	}
	return cmd
}

func doRunNewAccount() {
	//
	// Load up dependencies.
	//

	// ------ Common ------
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter()

	// ------ Repository ------
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	blockchainStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)
	tokRepo := repo.NewTokenRepo(cfg, logger, dbClient)
	gbdRepo := repo.NewGenesisBlockDataRepo(cfg, logger, dbClient)
	bdRepo := repo.NewBlockDataRepo(cfg, logger, dbClient)

	// ------ Use-case ------
	// Wallet
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	_ = walletEncryptKeyUseCase
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	_ = walletDecryptKeyUseCase
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)
	_ = createWalletUseCase

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)
	_ = createAccountUseCase
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		cfg,
		logger,
		accountRepo,
	)

	// Blockchain State
	getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)
	upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokRepo,
	)
	getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
		cfg,
		logger,
		tokRepo,
	)

	// Genesis BlockData
	upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
		cfg,
		logger,
		gbdRepo,
	)

	// BlockData
	upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
		cfg,
		logger,
		bdRepo,
	)

	// Proof of Work
	proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(
		cfg,
		logger,
	)

	// ------ Service ------
	getProofOfAuthorityPrivateKeyService := service.NewGetProofOfAuthorityPrivateKeyService(
		cfg,
		logger,
		getWalletUseCase,
		walletDecryptKeyUseCase,
	)
	createGenesisBlockDataService := service.NewCreateGenesisBlockDataService(
		cfg,
		logger,
		getProofOfAuthorityPrivateKeyService,
		getAccountUseCase,
		upsertAccountUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
		getAccountsHashStateUseCase,
		getTokensHashStateUseCase,
		proofOfWorkUseCase,
		upsertGenesisBlockDataUseCase,
		upsertBlockDataUseCase,
		upsertBlockchainStateUseCase,
		getBlockchainStateUseCase,
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
		// Execution
		blockchainState, err := createGenesisBlockDataService.Execute(sessCtx)
		if err != nil {
			logger.Error("Failed initializing new blockchain from new genesis block",
				slog.Any("error", err))
			return nil, err
		}
		if blockchainState == nil {
			err := errors.New("Blockchain state does not exist")
			return nil, err
		}

		return blockchainState, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed creating account: %v\n", err)
	}

	blockchainState := res.(*domain.BlockchainState)

	logger.Debug("Genesis block created",
		slog.Any("chain_id", blockchainState.ChainID),
	)
}
