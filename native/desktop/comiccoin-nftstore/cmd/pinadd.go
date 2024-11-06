package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	pkgfilepath "path/filepath"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	pkg_config "github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service"
)

// Command line argument flags
var (
	flagFilepath string
	// flagAPIKey   string
)

func PinAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pinadd",
		Short: "Commands used to upload a file to the NFT store and have it pinned as well",
		Run: func(cmd *cobra.Command, args []string) {
			doPinAddCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagFilepath, "filepath", "", "The path to the file you want to upload to the app")
	cmd.MarkFlagRequired("filepath")
	// cmd.Flags().StringVar(&flagAPIKey, "api-key", "", "The api-key to attach to the request")
	// cmd.MarkFlagRequired("api-key")

	return cmd
}

func doPinAddCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	// Environment variables.
	appSecretKey := config.GetEnvString("COMICCOIN_NFTSTORE_APP_SECRET_KEY", true)
	hmacSecretKey := config.GetEnvBytes("COMICCOIN_NFTSTORE_HMAC_SECRET_KEY", true)
	apiKey := config.GetEnvBytes("COMICCOIN_NFTSTORE_API_KEY", true)

	// --- Common --- //
	logger := logger.NewLogger()
	logger.Info("Starting file uploader ...",
		slog.Any("api_key", apiKey))

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
	_ = comicCoinConfig

	//
	// STEP 2:
	// Open the file.
	//

	// Get the filename from the filepath.
	fileName := pkgfilepath.Base(flagFilepath)

	// Open the file at the path.
	file, err := os.Open(flagFilepath)
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}
	defer file.Close()

	// Create a buffer to write the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	//
	// STEP 3
	// Add the file to the form.
	//

	// Create a form field writer for the file field
	fileField, err := writer.CreateFormFile("data", fileName)
	if err != nil {
		log.Fatalf("failed to create form file: %v", err)
	}

	// Copy the contents of the *os.File to the multipart form field
	if _, err := io.Copy(fileField, file); err != nil {
		log.Fatalf("failed to copy file to form field: %v", err)
	}

	//
	// STEP X
	// Close the form
	//

	// Close the multipart writer to finalize the form data
	if err := writer.Close(); err != nil {
		log.Fatalf("failed to close writer: %v", err)
	}

	contentType := writer.FormDataContentType()

	formData := &b

	// Send HTTP request with the multipart form data
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/ipfs/pin-add", formData)
	if err != nil {
		fmt.Printf("failed to create request: %v\n", err)
		return
	}

	// Create a Bearer string by appending string access token
	var bearer = "JWT " + string(apiKey)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Set the Content-Type header
	req.Header.Set("Content-Type", contentType)

	// disposition, params, err := mime.ParseMediaType(fmt.Sprintf("attachment;filename=%s", fileName))
	// fmt.Println("params", params)
	// fmt.Println("err", err)

	disposition := fmt.Sprintf("attachment;filename=%s", fileName)
	req.Header.Set("Content-Disposition", disposition)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

		// Try to decode the response as a string first
		var jsonStr string
		err := json.NewDecoder(teeReader).Decode(&jsonStr)
		if err != nil {
			logger.Error("decoding string error",
				slog.Any("err", err),
				slog.String("json", rawJSON.String()),
			)
			return
		}

		// Now try to decode the string into a map
		err = json.Unmarshal([]byte(jsonStr), &e)
		if err != nil {
			logger.Error("decoding map error",
				slog.Any("err", err),
				slog.String("json", jsonStr),
			)
			return
		}

		logger.Debug("Parsed error response",
			slog.Any("errors", e),
		)
		return
	}

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	post := &service.IPFSPinAddResponseIDO{}
	if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
		logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return
	}

	logger.Debug("Submitted successfully",
		slog.Any("RequestID", post.RequestID),
		slog.Any("Status", post.Status),
		slog.Any("Created", post.Created),
		slog.Any("Delegates", post.Delegates),
		slog.Any("Info", post.Info),
		slog.Any("CID", post.CID),
		slog.Any("Name", post.Name),
		slog.Any("Origins", post.Origins),
		slog.Any("Meta", post.Meta))
}