package blockchain

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

// Command line argument flags
var (
	flagIPFSIP               string
	flagIPFSPort             string
	flagIPFSPubGatewayDomain string
)

func DownloadTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "download",
		Short: "Download the token detail",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			logger := logger.NewLogger()
			logger.Debug("Excuting...",
				slog.String("data_dir", flagDataDir))
			if flagDataDir == "./data" {
				log.Fatal("cannot be `./data`")
			}

			cfg := &config.Config{
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
				IPFS: config.IPFSConfig{
					RemoteIP:            flagIPFSIP,
					RemotePort:          flagIPFSPort,
					PublicGatewayDomain: flagIPFSPubGatewayDomain,
				},
			}

			// --- Repositories ---

			nftokDB := disk.NewDiskStorage(flagDataDir, "non_fungible_token", logger)
			tokDB := disk.NewDiskStorage(flagDataDir, "token", logger)

			nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokDB)
			tokRepo := repo.NewTokenRepo(
				cfg,
				logger,
				tokDB)
			ipfsRepo := repo.NewIPFSRepo(cfg, logger)

			// --- Use-cases ---

			getTokUseCase := usecase.NewGetTokenUseCase(
				cfg,
				logger,
				tokRepo)

			getNFTokUseCase := usecase.NewGetNonFungibleTokenUseCase(
				cfg,
				logger,
				nftokenRepo)

			downloadNFTokMetadataUsecase := usecase.NewDownloadMetadataNonFungibleTokenUseCase(
				cfg,
				logger,
				ipfsRepo)

			downloadNFTokAssetUsecase := usecase.NewDownloadNonFungibleTokenAssetUseCase(
				cfg,
				logger,
				ipfsRepo)

			upsertNFTokUseCase := usecase.NewUpsertNonFungibleTokenUseCase(
				cfg,
				logger,
				nftokenRepo)

			// --- Service ---

			getOrDownloadNonFungibleTokenService := service.NewGetOrDownloadNonFungibleTokenService(
				cfg,
				logger,
				getNFTokUseCase,
				getTokUseCase,
				downloadNFTokMetadataUsecase,
				downloadNFTokAssetUsecase,
				upsertNFTokUseCase)

			//
			// STEP 2
			// Check if we can connect with IPFS node.
			//

			peerID, err := ipfsRepo.ID()
			if err != nil {
				log.Fatalf("Failed connecting to IPFS node, you are not connected.")
			}
			fmt.Printf("IPFS Node ID: %s\n", peerID)

			//
			// STEP 3
			// Lookup our `token id` in our NFT db and if it exists we can
			// exit this command as we've already downloaded the data.
			//

			tokenID, err := strconv.ParseUint(flagTokenID, 10, 64)
			if err != nil {
				log.Fatalf("failed converting token id to unit64: %v\n", err)
			}

			nftok, err := getOrDownloadNonFungibleTokenService.Execute(tokenID)
			if err != nil {
				log.Fatalf("Failed downloading non-fungible tokens: %v\n", err)
			}

			logger.Debug("Downloaded NFT successfully.",
				slog.Any("token_id", nftok.TokenID),
				slog.Any("metadata_uri", nftok.MetadataURI),
				slog.Any("metadata", nftok.Metadata),
			)
		},
	}
	cmd.Flags().StringVar(&flagDataDir, "datadir", config.GetDefaultDataDirectory(), "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The value to lookup the token by")
	cmd.MarkFlagRequired("token-id")
	cmd.Flags().StringVar(&flagIPFSIP, "ipfs-ip", "127.0.0.1", "")
	cmd.Flags().StringVar(&flagIPFSPort, "ipfs-port", "5001", "")
	cmd.Flags().StringVar(&flagIPFSPubGatewayDomain, "ipfs-public-gateway-domain", "https://ipfs.io", "")
	return cmd
}
