package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"path/filepath"
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

			nftokenByTokenIDDB := disk.NewDiskStorage(flagDataDir, "non_fungible_token_by_id", logger)
			nftokenByMetadataURIDB := disk.NewDiskStorage(flagDataDir, "non_fungible_token_by_metadata_uri", logger)
			tokDB := disk.NewDiskStorage(flagDataDir, "token", logger)

			nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokenByTokenIDDB, nftokenByMetadataURIDB)
			tokRepo := repo.NewTokenRepo(
				cfg,
				logger,
				tokDB)
			ipfsNode := repo.NewIPFSRepo(cfg, logger)

			getTokUseCase := usecase.NewGetTokenUseCase(
				cfg,
				logger,
				tokRepo)

			// //
			// // STEP 2
			// // Check if we can connect with IPFS node.
			// //
			//
			// peerID, err := ipfsNode.ID()
			// if err != nil {
			// 	log.Fatalf("failed connecting to IPFS repo to get ID(): %v\n", err)
			// }
			// fmt.Printf("IPFS Node ID: %s\n", peerID)

			//
			// STEP 3
			// Lookup our `token id` in our NFT db and if it exists we can
			// exit this command as we've already downloaded the data.
			//

			tokenID, err := strconv.ParseUint(flagTokenID, 10, 64)
			if err != nil {
				log.Fatalf("failed converting token id to unit64: %v\n", err)
			}

			nftok, err := nftokenRepo.GetByTokenID(tokenID)
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

			//
			// STEP 5
			// Download the metadata file and unmarshal it.
			//

			cid := strings.Replace(metadataURI, "ipfs://", "", -1)
			metadataBytes, contentType, err := ipfsNode.Get(context.Background(), cid)
			if err != nil {
				log.Fatalf("failed getting from cid: %v\n", err)
			}

			var metadata *domain.NonFungibleTokenMetadata
			if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
				log.Fatalf("failed unmarshalling metadata: %v\n", err)
			}

			logger.Debug("Downloaded metadata",
				slog.Any("cid", cid),
				slog.Any("metadata", metadata),
				slog.Any("content_type", contentType))

			//
			// STEP 6
			// Save the metadata in local filesystem
			//

			metadataFilepath := filepath.Join(flagDataDir, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), "metadata.json")

			// Create the directories recursively.
			if err := os.MkdirAll(filepath.Dir(metadataFilepath), 0755); err != nil {
				log.Fatalf("Failed create directories: %v\n", err)
			}

			if err := ioutil.WriteFile(metadataFilepath, metadataBytes, 0644); err != nil {
				log.Fatalf("Failed write metadata file: %v\n", err)
			}

			//
			// STEP 7
			// Download the image file from IPFS and save locally.
			//

			imageCID := strings.Replace(metadata.Image, "ipfs://", "", -1)
			imageBytes, imageContentType, err := ipfsNode.Get(context.Background(), imageCID)
			if err != nil {
				log.Fatalf("failed getting image via cid: %v\n", err)
			}

			var imageFilename string
			switch imageContentType {
			case "image/png": //  Portable Network Graphics (PNG)
				imageFilename = fmt.Sprintf("%v.png", imageCID)
			case "image/jpeg": // Joint Photographic Experts Group (JPEG)
				imageFilename = fmt.Sprintf("%v.jpg", imageCID)
			case "image/jpg": // Joint Photographic Experts Group (JPG)
				imageFilename = fmt.Sprintf("%v.jpg", imageCID)
			case "image/gif": //  Graphics Interchange Format (GIF)
				imageFilename = fmt.Sprintf("%v.gif", imageCID)
			case "image/bmp": // Bitmap (BMP)
				imageFilename = fmt.Sprintf("%v.bmp", imageCID)
			case "image/tiff": // Tagged Image File Format (TIFF)
				imageFilename = fmt.Sprintf("%v.tiff", imageCID)
			case "image/tif": // Tagged Image File Format (TIFF)
				imageFilename = fmt.Sprintf("%v.tif", imageCID)
			case "image/webp": // Web Picture (WEBP)
				imageFilename = fmt.Sprintf("%v.webp", imageCID)
			case "image/svg+xml": // Scalable Vector Graphics (SVG)
				imageFilename = fmt.Sprintf("%v.svg", imageCID)
			default:
				log.Fatalf("Unsupported image type: %v\n", imageContentType)
			}

			imageFilepath := filepath.Join(flagDataDir, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), imageFilename)
			// Save the data to file.
			f, err := os.Create(imageFilepath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			// Convert the response bytes into reader.
			imageBytesReader := bytes.NewReader(imageBytes)

			// Save to local directory.
			_, err = io.Copy(f, imageBytesReader)
			if err != nil {
				log.Fatal(err)
			}

			//
			// STEP 8
			// Download the animation file from IPFS and save locally.
			//

			animationCID := strings.Replace(metadata.AnimationURL, "ipfs://", "", -1)
			animationBytes, animationContentType, err := ipfsNode.Get(context.Background(), animationCID)
			if err != nil {
				log.Fatalf("failed getting animation via cid: %v\n", err)
			}

			var animationFilename string
			switch animationContentType {
			case "video/mp4": // MPEG-4
				animationFilename = fmt.Sprintf("%v.mp4", animationCID)
			case "video/x-m4v": // MPEG-4 Video
				animationFilename = fmt.Sprintf("%v.m4v", animationCID)
			case "video/quicktime": // QuickTime
				animationFilename = fmt.Sprintf("%v.mov", animationCID)
			case "video/webm": // WebM
				animationFilename = fmt.Sprintf("%v.webm", animationCID)
			case "video/ogg": // Ogg Theora
				animationFilename = fmt.Sprintf("%v.ogv", animationCID)
			case "video/x-flv": // Flash Video
				animationFilename = fmt.Sprintf("%v.flv", animationCID)
			case "video/x-msvideo": // AVI (Audio Video Interleave)
				animationFilename = fmt.Sprintf("%v.avi", animationCID)
			case "video/x-ms-wmv": //  Windows Media Video
				animationFilename = fmt.Sprintf("%v.wmv", animationCID)
			case "video/3gpp": // 3GPP (3rd Generation Partnership Project)
				animationFilename = fmt.Sprintf("%v.3gp", animationCID)
			case "video/3gpp2": // 3GPP2 (3rd Generation Partnership Project 2)
				animationFilename = fmt.Sprintf("%v.3g2", animationCID)
			default:
				log.Fatalf("Unsupported video type: %v\n", animationContentType)
			}

			animationFilepath := filepath.Join(flagDataDir, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), animationFilename)
			// Save the data to file.
			f, err = os.Create(animationFilepath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			// Convert the response bytes into reader.
			animationBytesReader := bytes.NewReader(animationBytes)

			// Save to local directory.
			_, err = io.Copy(f, animationBytesReader)
			if err != nil {
				log.Fatal(err)
			}

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

			// if err := nftokenRepo.Upsert(nftok); err != nil {
			// 	log.Fatalf("Failed upserting nft token: %v\n", err)
			// }

		},
	}
	cmd.Flags().StringVar(&flagDataDir, "datadir", config.GetDefaultDataDirectory(), "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The value to lookup the token by")
	cmd.MarkFlagRequired("token-id")
	cmd.Flags().StringVar(&flagIPFSIP, "ipfs-ip", "", "")
	cmd.MarkFlagRequired("ipfs-ip")
	cmd.Flags().StringVar(&flagIPFSPort, "ipfs-port", "", "")
	cmd.MarkFlagRequired("ipfs-port")
	cmd.Flags().StringVar(&flagIPFSPubGatewayDomain, "ipfs-public-gateway-domain", "https://ipfs.io", "")
	return cmd
}
