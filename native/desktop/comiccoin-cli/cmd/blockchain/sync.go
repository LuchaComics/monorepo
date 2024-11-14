package blockchain

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

// Command line argument flags
var (
	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

func BlockchainSyncCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Execute command to synchronize the local blockchain with the Blockchain network.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doRunBlockchainSyncCmd(); err != nil {
				log.Fatalf("Failed to sync blockchain: %v\n", err)
			}
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	return cmd
}

func doRunBlockchainSyncCmd() error {
	// ------ Common ------
	logger := logger.NewProvider()
	logger.Info("Syncing blockchain...",
		slog.Any("authority_address", flagAuthorityAddress))

	// // ------ Database -----
	// // walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	// blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	// blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	//
	// // ------------ Repo ------------
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		logger,
		genesisBlockDataDB)
	// blockchainStateRepo := repo.NewBlockchainStateRepo(
	// 	cfg,
	// 	logger,
	// 	blockchainStateDB)
	// blockchainStateDTORepo := repo.NewBlockchainStateDTORepo(
	// 	cfg,
	// 	logger)

	genesisBlockDataDTORepoConfig := auth_repo.NewGenesisBlockDataDTOConfigurationProvider(flagAuthorityAddress)
	genesisBlockDataDTORepo := auth_repo.NewGenesisBlockDataDTORepo(
		genesisBlockDataDTORepoConfig,
		logger)

	// blockDataRepo := repo.NewBlockDataRepo(
	// 	cfg,
	// 	logger,
	// 	blockDataDB)
	// blockDataDTORepo := repo.NewBlockDataDTORepo(
	// 	cfg,
	// 	logger)
	//
	// _ = blockDataRepo
	//
	// // ------------ Use-Case ------------
	// // Blockchain State
	// upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
	// 	cfg,
	// 	logger,
	// 	blockchainStateRepo)
	// getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
	// 	cfg,
	// 	logger,
	// 	blockchainStateRepo)
	//
	// // Blockchain State DTO
	// getBlockchainStateFromCentralAuthorityByChainIDUseCase := usecase.NewGetBlockchainStateFromCentralAuthorityByChainIDUseCase(
	// 	cfg,
	// 	logger,
	// 	blockchainStateDTORepo)

	// Genesis Block Data DTO
	getGenesisBlockDataDTOFromCentralAuthorityUseCase := auth_usecase.NewGetGenesisBlockDataDTOFromBlockchainAuthorityUseCase(
		logger,
		genesisBlockDataDTORepo)

	// // Genesis Block Data
	upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
		logger,
		genesisBlockDataRepo)
	getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
		logger,
		genesisBlockDataRepo)

	// // Block Data DTO
	// getBlockDataFromCentralAuthorityByBlockNumberUseCase := usecase.NewGetBlockDataFromCentralAuthorityByBlockNumberUseCase(
	// 	cfg,
	// 	logger,
	// 	blockDataDTORepo)
	//
	// ------------ Service ------------
	genesisBlockDataGetService := service.NewGenesisBlockDataGetService(
		logger,
		getGenesisBlockDataUseCase,
		upsertGenesisBlockDataUseCase,
		getGenesisBlockDataDTOFromCentralAuthorityUseCase,
	)

	blockchainSyncService := service.NewBlockchainSyncWithBlockchainAuthorityService(
		logger,
		genesisBlockDataGetService,
	)

	// ------------ Execute ------------

	ctx := context.Background()

	if err := blockchainSyncService.Execute(ctx, flagChainID); err != nil {
		log.Fatalf("Failed to sync blockchain: %v\n", err)
	}

	logger.Info("Finished syncing blockchain",
		slog.Any("authority_address", flagAuthorityAddress))
	return nil
}
