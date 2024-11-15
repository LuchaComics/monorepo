package coins

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

// Command line argument flags
var (
	flagKeystoreFile                  string // Location of the wallet keystore
	flagDataDir                       string // Location of the database directory
	flagLabel                         string
	flagSenderAccountAddress          string
	flagSenderAccountPassword         string
	flagSenderAccountPasswordRepeated string
	flagCoinbaseAddress               string
	flagRecipientAddress              string
	flagQuantity                      uint64
	flagKeypairName                   string
	flagData                          string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int

	flagListenHTTPAddress string

	flagIdentityKeyID string

	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

func TransferCoinsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "transfer",
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network to transfer coins from coinbase account to another account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunTransferCoinsCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagSenderAccountAddress, "sender-account-address", "", "The address of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("sender-account-address")

	cmd.Flags().StringVar(&flagSenderAccountPassword, "sender-account-password", "", "The password to unlock the account which will transfer the coin")
	cmd.MarkFlagRequired("sender-account-password")

	cmd.Flags().Uint64Var(&flagQuantity, "value", 0, "The amount of coins to send")
	cmd.MarkFlagRequired("value")

	cmd.Flags().StringVar(&flagData, "data", "", "Optional data to include with this transaction")

	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The address of the account whom will receive this coin")
	cmd.MarkFlagRequired("recipient-address")

	return cmd
}

func doRunTransferCoinsCommand() {
	logger := logger.NewProvider()

	// ------ Common ------
	keystore := keystore.NewAdapter()
	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	tokenRepo := disk.NewDiskStorage(flagDataDirectory, "token", logger)

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
		tokenRepo)
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

	_ = walletEncryptKeyUseCase
	_ = createWalletUseCase
	_ = createWalletUseCase
	_ = listAllWalletUseCase
	_ = createAccountUseCase
	_ = getAccountsHashStateUseCase
	_ = upsertAccountUseCase

	// ------ Service ------

	coinTransferService := service.NewCoinTransferService(
		logger,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
	)

	// ------ Execute ------
	ctx := context.Background()
	recAddr := common.HexToAddress(strings.ToLower(flagRecipientAddress))
	sendAddr := common.HexToAddress(strings.ToLower(flagSenderAccountAddress))

	if err := storageTransactionOpenUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	coinTransferServiceErr := coinTransferService.Execute(
		ctx,
		flagChainID,
		&sendAddr,
		flagSenderAccountPassword,
		&recAddr,
		flagQuantity, // A.k.a. `value`.
		[]byte(flagData),
	)
	if coinTransferServiceErr != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed execute coin transfer service: %v", coinTransferServiceErr)
	}

	if err := storageTransactionCommitUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	logger.Info("Finished transfering coin(s) to another account",
		slog.Any("data-director", flagDataDirectory),
		slog.Any("chain-id", flagChainID),
		slog.Any("nftstorage-address", flagNFTStorageAddress),
		slog.Any("sender-account-address", flagSenderAccountAddress),
		slog.Any("sender-account-password", flagSenderAccountPassword),
		slog.Any("value", flagQuantity),
		slog.Any("data", flagData),
		slog.Any("recipient-address", flagRecipientAddress),
		slog.Any("authority-address", flagAuthorityAddress))
}
