package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/interface/http/handler"
)

// Command line argument flags
var ()

func RemoteVersionCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remoteversion",
		Short: "Commands fetches version from api",
		Run: func(cmd *cobra.Command, args []string) {
			doRemoteVersionCmd()
		},
	}
	return cmd
}

const (
	versionURL = "/version"
)

func doRemoteVersionCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	listenHTTPAddress := config.GetEnvString("COMICCOIN_NFTSTORAGE_ADDRESS", true)

	logger := logger.NewProvider()
	httpEndpoint := fmt.Sprintf("http://%s%s", listenHTTPAddress, versionURL)

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
	resp, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do get request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Fatalf("http endpoint does not exist for: %v", httpEndpoint)
	}

	if resp.StatusCode == http.StatusBadRequest {
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

	//
	// STEP 3:
	// Print the success message.
	//

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	respContent := &httphandler.VersionResponseIDO{}
	if err := json.NewDecoder(teeReader).Decode(&respContent); err != nil {
		logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return
	}

	logger.Debug("Request successfully received",
		slog.Any("response", respContent))
}
