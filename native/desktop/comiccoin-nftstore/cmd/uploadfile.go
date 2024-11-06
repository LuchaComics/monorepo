package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
	pkgfilepath "path/filepath"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	pkg_config "github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	pkg_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/usecase"
)

// Command line argument flags
var (
	flagFilepath string
)

func UploadFileCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "uploadfile",
		Short: "Commands used to upload a file to the NFT store",
		Run: func(cmd *cobra.Command, args []string) {
			doUploadFileCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagFilepath, "filepath", "", "The path to the file you want to upload to the app")
	cmd.MarkFlagRequired("filepath")

	return cmd
}

func doUploadFileCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	// Environment variables.
	appSecretKey := config.GetEnvString("COMICCOIN_NFTSTORE_APP_SECRET_KEY", true)
	hmacSecretKey := config.GetEnvBytes("COMICCOIN_NFTSTORE_HMAC_SECRET_KEY", true)

	// --- Common --- //
	logger := logger.NewLogger()
	logger.Info("Starting file uploader ...",
		slog.Any("flatHMACSecret", hmacSecretKey),
		slog.Any("appSecretKey", appSecretKey))

	comicCoinConfig := &pkg_config.Config{
		IPFS: pkg_config.IPFSConfig{
			RemoteIP:            constants.ComicCoinIPFSRemoteIP,
			RemotePort:          constants.ComicCoinIPFSRemotePort,
			PublicGatewayDomain: constants.ComicCoinIPFSPublicGatewayDomain,
		},
	}

	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                        constants.ComicCoinChainID,
			TransPerBlock:                  constants.ComicCoinTransPerBlock,
			Difficulty:                     constants.ComicCoinDifficulty,
			ConsensusPollingDelayInMinutes: constants.ComicCoinConsensusPollingDelayInMinutes,
			ConsensusProtocol:              constants.ComicCoinConsensusProtocol,
		},
		App: config.AppConfig{
			DirPath:     flagDataDir,
			HTTPAddress: flagListenHTTPAddress,
			HMACSecret:  hmacSecretKey,
			AppSecret:   appSecretKey,
		},
		DB: config.DBConfig{
			DataDir: flagDataDir,
		},
		Peer: config.PeerConfig{
			ListenPort: constants.ComicCoinPeerListenPort,
			KeyName:    constants.ComicCoinIdentityKeyID,
		},
		IPFS: config.IPFSConfig{
			RemoteIP:            constants.ComicCoinIPFSRemoteIP,
			RemotePort:          constants.ComicCoinIPFSRemotePort,
			PublicGatewayDomain: constants.ComicCoinIPFSPublicGatewayDomain,
		},
	}

	_ = cfg

	// --- Repository --- //

	ipfsRepo := pkg_repo.NewIPFSRepo(comicCoinConfig, logger)

	// --- UseCase --- //

	ipfsPinAddUsecase := usecase.NewIPFSPinAddUseCase(logger, ipfsRepo)

	// --- Service --- //

	createPinObjectService := service.NewCreatePinObjectService(
		logger,
		ipfsPinAddUsecase,
	)
	_ = createPinObjectService

	// --- Execute our command functionality --- //
	// Open the file at the path.
	file, err := os.Open(flagFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read the response body as a byte slice
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the filename from the filepath.
	filename := pkgfilepath.Base(flagFilepath)

	// Get the content type of the file.
	contentType := http.DetectContentType(bytes)

	// Populate the required `meta` for IPFS.
	meta := make(map[string]string, 0)
	meta["filename"] = filename
	meta["content_type"] = contentType

	req := &service.PinObjectCreateServiceRequestIDO{
		Name:    filename,
		Origins: make([]string, 0),
		Meta:    meta,
		File:    file,
	}

	logger.Debug("Uploading now...", slog.Any("req", req))

	res, err := createPinObjectService.Execute(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed submitting file to pin in ipfs: %v\n", err)
	}

	logger.Debug("Successfully uploaded file",
		slog.Any("res", res))
}
