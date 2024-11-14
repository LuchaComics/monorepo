package blockchain

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"
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
	// genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	// blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	// blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	//
	// // ------------ Repo ------------
	// genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
	// 	cfg,
	// 	logger,
	// 	genesisBlockDataDB)
	// blockchainStateRepo := repo.NewBlockchainStateRepo(
	// 	cfg,
	// 	logger,
	// 	blockchainStateDB)
	// blockchainStateDTORepo := repo.NewBlockchainStateDTORepo(
	// 	cfg,
	// 	logger)
	// genesisBlockDataDTORepo := repo.NewGenesisBlockDataDTORepo(
	// 	cfg,
	// 	logger)
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
	//
	// // Genesis Block Data DTO
	// getGenesisBlockDataFromCentralAuthorityByChainIDUseCase := usecase.NewGetGenesisBlockDataFromCentralAuthorityByChainIDUseCase(
	// 	cfg,
	// 	logger,
	// 	genesisBlockDataDTORepo)
	//
	// // Genesis Block Data
	// upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	genesisBlockDataRepo)
	// getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	genesisBlockDataRepo)
	//
	// // Block Data DTO
	// getBlockDataFromCentralAuthorityByBlockNumberUseCase := usecase.NewGetBlockDataFromCentralAuthorityByBlockNumberUseCase(
	// 	cfg,
	// 	logger,
	// 	blockDataDTORepo)
	//
	// // ------------ Service ------------
	// localBlockchainSyncService := service.NewLocalBlockchainSyncWithCentralAuthorityService(
	// 	cfg,
	// 	logger,
	// 	getBlockchainStateUseCase,
	// 	getBlockchainStateFromCentralAuthorityByChainIDUseCase,
	// 	upsertBlockchainStateUseCase,
	// 	getGenesisBlockDataUseCase,
	// 	getGenesisBlockDataFromCentralAuthorityByChainIDUseCase,
	// 	upsertGenesisBlockDataUseCase,
	// 	getBlockDataFromCentralAuthorityByBlockNumberUseCase,
	// )
	//
	// ctx := context.Background()
	//
	// if err := localBlockchainSyncService.Execute(ctx); err != nil {
	// 	log.Fatalf("Failed to sync local blockchain: %v\n", err)
	// }

	// ------------ Execute ------------
	logger.Info("Finished syncing blockchain",
		slog.Any("authority_address", flagAuthorityAddress))
	return nil
}
