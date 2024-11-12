package blockchain

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

// Command line argument flags
var (
	flagDataDirectory string
	flagChainID       string
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

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", "", "The data directory to save to")
	cmd.MarkFlagRequired("data-directory")

	return cmd
}

func doRunBlockchainSyncCmd() error {
	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                     constants.ChainIDMainNet,
			CentralAuthorityHTTPAddress: "http://127.0.0.1:8000",
		},
		App: config.AppConfig{
			DirPath:     flagDataDirectory,
			HTTPAddress: "http://127.0.0.1:8000",
		},
	}

	// ------ Common ------
	logger := logger.NewProvider()

	// ------ Database -----
	// walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)

	// ------------ Repo ------------
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		cfg,
		logger,
		genesisBlockDataDB)
	blockDataRepo := repo.NewBlockDataRepo(
		cfg,
		logger,
		blockDataDB)
	blockchainStateRepo := repo.NewBlockchainStateRepo(
		cfg,
		logger,
		blockchainStateDB)
	blockchainStateDTORepo := repo.NewBlockchainStateDTORepo(
		cfg,
		logger)
	genesisBlockDataDTORepo := repo.NewGenesisBlockDataDTORepo(
		cfg,
		logger)

	// ------------ Use-Case ------------
	// Blockchain State DTO
	getBlockchainStateFromCentralAuthorityByChainIDUseCase := usecase.NewGetBlockchainStateFromCentralAuthorityByChainIDUseCase(
		cfg,
		logger,
		blockchainStateDTORepo)

	// Genesis Block Data DTO
	getGenesisBlockDataFromCentralAuthorityByChainIDUseCase := usecase.NewGetGenesisBlockDataFromCentralAuthorityByChainIDUseCase(
		cfg,
		logger,
		genesisBlockDataDTORepo)

	// Genesis Block Data
	upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
		cfg,
		logger,
		genesisBlockDataRepo)
	getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
		cfg,
		logger,
		genesisBlockDataRepo)

	_ = blockDataRepo
	_ = blockchainStateRepo

	// ------------ Service ------------
	localBlockchainSyncService := service.NewLocalBlockchainSyncService(
		cfg,
		logger,
		getBlockchainStateFromCentralAuthorityByChainIDUseCase,
		getGenesisBlockDataUseCase,
		getGenesisBlockDataFromCentralAuthorityByChainIDUseCase,
		upsertGenesisBlockDataUseCase,
	)

	ctx := context.Background()

	if err := localBlockchainSyncService.Execute(ctx); err != nil {
		log.Fatalf("Failed to sync local blockchain: %v\n", err)
	}

	// ------------ Execute ------------
	return nil
}
