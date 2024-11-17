package tokens

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

func BurnTokensCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "burn",
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network to transfer tokens from your account to another account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunBurnTokensCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagSenderAccountAddress, "sender-account-address", "", "The address of the account we will use in our token transfer")
	cmd.MarkFlagRequired("sender-account-address")

	cmd.Flags().StringVar(&flagSenderAccountPassword, "sender-account-password", "", "The password to unlock the account which will transfer the token")
	cmd.MarkFlagRequired("sender-account-password")

	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The unique token identification to use to lookup the token")
	cmd.MarkFlagRequired("token-id")

	return cmd
}

func doRunBurnTokensCommand() {
	// ------ Common ------
	logger := logger.NewProvider()
	keystore := keystore.NewAdapter()
	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	tokenDB := disk.NewDiskStorage(flagDataDirectory, "token", logger)

	// ------ Repo ------
	walletRepo := repo.NewWalletRepo(
		logger,
		walletDB)
	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		logger,
		genesisBlockDataDB)
	blockchainStateRepo := repo.NewBlockchainStateRepo(
		logger,
		blockchainStateDB)
	blockDataRepo := repo.NewBlockDataRepo(
		logger,
		blockDataDB)
	tokRepo := repo.NewTokenRepo(
		logger,
		tokenDB)
	mempoolTxDTORepoConfig := auth_repo.NewMempoolTransactionDTOConfigurationProvider(flagAuthorityAddress)
	mempoolTxDTORepo := auth_repo.NewMempoolTransactionDTORepo(mempoolTxDTORepoConfig, logger)

	// ------ Use-case ------

	// Storage Transaction
	storageTransactionOpenUseCase := usecase.NewStorageTransactionOpenUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo)
	storageTransactionCommitUseCase := usecase.NewStorageTransactionCommitUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo)
	storageTransactionDiscardUseCase := usecase.NewStorageTransactionDiscardUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo)

	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		logger,
		walletRepo)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		logger,
		walletRepo)
	listAllWalletUseCase := usecase.NewListAllWalletUseCase(
		logger,
		walletRepo)

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		logger,
		accountRepo)

	// Mempool Transaction DTO
	submitMempoolTransactionDTOToBlockchainAuthorityUseCase := auth_usecase.NewSubmitMempoolTransactionDTOToBlockchainAuthorityUseCase(
		logger,
		mempoolTxDTORepo,
	)

	// Token
	getTokenUseCase := usecase.NewGetTokenUseCase(
		logger,
		tokRepo,
	)

	_ = walletEncryptKeyUseCase
	_ = createWalletUseCase
	_ = createWalletUseCase
	_ = listAllWalletUseCase
	_ = createAccountUseCase
	_ = getAccountsHashStateUseCase
	_ = upsertAccountUseCase

	// ------ Service ------

	tokenBurnService := service.NewTokenBurnService(
		logger,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getTokenUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
	)

	// ------ Execute ------
	ctx := context.Background()
	sendAddr := common.HexToAddress(strings.ToLower(flagSenderAccountAddress))
	tokenID, ok := new(big.Int).SetString(flagTokenID, 10)
	if !ok {
		log.Fatal("Failed convert `token_id` to big.Int")
	}

	logger.Debug("Transfering Token...",
		slog.Any("token_id", tokenID))

	if err := storageTransactionOpenUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	tokenBurnServiceErr := tokenBurnService.Execute(
		ctx,
		flagChainID,
		&sendAddr,
		flagSenderAccountPassword,
		tokenID,
	)
	if tokenBurnServiceErr != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed execute token transfer service: %v", tokenBurnServiceErr)
	}

	if err := storageTransactionCommitUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	logger.Info("Finished burning token",
		slog.Any("data-director", flagDataDirectory),
		slog.Any("chain-id", flagChainID),
		slog.Any("nftstorage-address", flagNFTStorageAddress),
		slog.Any("sender-account-address", flagSenderAccountAddress),
		slog.Any("sender-account-password", flagSenderAccountPassword),
		slog.Any("token-id", flagTokenID),
		slog.Any("recipient-address", flagRecipientAddress),
		slog.Any("authority-address", flagAuthorityAddress))
}
