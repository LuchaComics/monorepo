package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
)

// Command line argument flags
var (
	flagCID string
)

func GetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Commands used get the contents of a file via the CID",
		Run: func(cmd *cobra.Command, args []string) {
			doGetCmd()
		},
	}

	cmd.Flags().StringVar(&flagCID, "cid", "", "The unique content ID of the file on the IPFS network")
	cmd.MarkFlagRequired("cid")

	return cmd
}

const (
	ipfsGatewayURL = "/ipfs/${CID}"
)

func doGetCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	listenHTTPAddress := config.GetEnvString("COMICCOIN_NFTSTORAGE_ADDRESS", true)

	logger := logger.NewProvider()
	modifiedIpfsGatewayURL := strings.ReplaceAll(ipfsGatewayURL, "${CID}", flagCID)
	httpEndpoint := fmt.Sprintf("http://%s%s", listenHTTPAddress, modifiedIpfsGatewayURL)

	//
	// STEP 2
	//

	r, err := http.NewRequest("GET", httpEndpoint, nil)
	if err != nil {
		log.Fatalf("failed to setup get request: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")

	logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "GET"))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		log.Fatalf("http endpoint does not exist for: %v", httpEndpoint)
	}

	if res.StatusCode == http.StatusBadRequest {
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

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

	//
	// STEP 3:
	// Read and process the content received from the server
	//

	// Determine the content type for conditional processing
	contentType := res.Header.Get("Content-Type")

	// If the content type is not JSON, assume it's file data and handle accordingly
	if contentType != "application/json" {
		// Save the output to a file or display it
		var outputFileName = "output-file" // default filename
		contentDisposition := res.Header.Get("Content-Disposition")
		if strings.Contains(contentDisposition, "filename*=") {
			// Extract filename from Content-Disposition header if available
			filenameParts := strings.Split(contentDisposition, "filename*=")
			if len(filenameParts) > 1 {
				outputFileName = strings.Trim(filenameParts[1], "\"")
			}
		}

		// Create the output file
		file, err := os.Create(outputFileName)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer file.Close()

		// Copy the response body to the file
		_, err = io.Copy(file, res.Body)
		if err != nil {
			log.Fatalf("failed to save file content: %v", err)
		}

		log.Printf("File saved as %s", outputFileName)
	} else {
		// If it's JSON or an error message, decode it and log the output
		var jsonResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
			logger.Error("failed to decode JSON response",
				slog.Any("error", err),
			)
		} else {
			logger.Debug("Received JSON response",
				slog.Any("data", jsonResponse),
			)
		}
	}
}
