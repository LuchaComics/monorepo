package blockchain

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
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
					LocalIP:             flagIPFSIP,
					LocalPort:           flagIPFSPort,
					PublicGatewayDomain: flagIPFSPubGatewayDomain,
				},
			}

			// --- Repositories ---

			nftokenByTokenIDDB := disk.NewDiskStorage(flagDataDir, "non_fungible_token_by_id", logger)
			nftokenByMetadataURIDB := disk.NewDiskStorage(flagDataDir, "non_fungible_token_by_metadata_uri", logger)
			tokDB := disk.NewDiskStorage(flagDataDir, "token", logger)

			nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokenByTokenIDDB, nftokenByMetadataURIDB)
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

			createNFTokUseCase := usecase.NewCreateNonFungibleTokenUseCase(
				cfg,
				logger,
				nftokenRepo)

			//
			// STEP 2
			// Check if we can connect with IPFS node.
			//

			peerID, err := ipfsRepo.ID()
			if err != nil {
				log.Fatalf("failed connecting to IPFS repo to get ID(): %v\n", err)
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

			nftok, err := getNFTokUseCase.Execute(tokenID)
			if err != nil {
				logger.Debug("err", slog.Any("error", err))
				return
			}
			if nftok != nil {
				logger.Debug("Token already exists locally, aborting...")
				return
			}

			//
			// STEP 4
			// Lookup our `token` in our db and retrieve the record so we can
			// extract the `Metadata URI` value necessary to lookup later in
			// the decentralized storage service (IPFS).
			//

			tok, err := getTokUseCase.Execute(tokenID)
			if err != nil {
				log.Fatalf("failed getting token due to err: %v\n", err)
			}
			if tok == nil {
				log.Fatalf("Token does not exist for: %v\n", tokenID)
			}

			metadataURI := tok.MetadataURI

			// Confirm URI is using protocol our app supports.
			if strings.Contains(metadataURI, "ipfs://") == false {
				log.Fatalf("Token metadata URI contains protocol we do not support: %v\n", metadataURI)
			}

			metadata, metadataFilepath, err := downloadNFTokMetadataUsecase.Execute(tok.ID, metadataURI)
			if err != nil {
				log.Fatalf("failed getting or downloading nft metadata: %v\n", err)
			}

			// Replace the IPFS path with our local systems filepath.
			metadataURI = metadataFilepath

			//
			// STEP 7
			// Download the image file from IPFS and save locally.
			//

			imageCID := strings.Replace(metadata.Image, "ipfs://", "", -1)
			imageFilepath, err := downloadNFTokAssetUsecase.Execute(tok.ID, imageCID)
			if err != nil {
				log.Fatalf("failed getting or downloading nft image asset: %v\n", err)
			}

			// Replace the IPFS path with our local systems filepath.
			metadata.Image = imageFilepath

			//
			// STEP 8
			// Download the animation file from IPFS and save locally.
			//

			animationCID := strings.Replace(metadata.AnimationURL, "ipfs://", "", -1)
			animationFilepath, err := downloadNFTokAssetUsecase.Execute(tok.ID, animationCID)
			if err != nil {
				log.Fatalf("failed getting or downloading nft image asset: %v\n", err)
			}

			// Replace the IPFS path with our local systems filepath.
			metadata.AnimationURL = animationFilepath

			//
			// STEP 8
			// Create our NFT token to be referenced in future.
			//

			nftok = &domain.NonFungibleToken{
				TokenID:     tokenID,
				MetadataURI: metadataURI,
				Metadata:    metadata,
			}

			logger.Debug("Downloaded from NFT.",
				slog.Any("token_id", nftok.TokenID),
				slog.Any("metadata_uri", nftok.MetadataURI),
				slog.Any("metadata", nftok.Metadata),
			)

			if err := createNFTokUseCase.Execute(nftok); err != nil {
				log.Fatalf("Failed creating nft token: %v\n", err)
			}
			// _ = createNFTokUseCase
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
