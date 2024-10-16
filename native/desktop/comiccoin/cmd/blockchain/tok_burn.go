package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/logger"
	handler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/interface/http/handler"
)

// Command line argument flags
var (
	flagBurnTokenOwnerAddress  string
	flagBurnTokenOwnerPassword string
	flagBurnTokenID            uint64
)

func httpJsonApiBurnTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "burn",
		Short: "(PoA only) Burn your non-fungible token",
		Run: func(cmd *cobra.Command, args []string) {
			doBurnToken()
		},
	}

	cmd.Flags().StringVar(&flagBurnTokenOwnerAddress, "token-owner-address", "", "(Required for `PoA` consensus protocol) The address of the token owner's account")
	cmd.MarkFlagRequired("token-owner-address")
	cmd.Flags().StringVar(&flagBurnTokenOwnerPassword, "token-owner-password", "", "(Required for `PoA` consensus protocol) The password in the token owner's wallet")
	cmd.MarkFlagRequired("token-owner-password")
	cmd.Flags().Uint64Var(&flagBurnTokenID, "token-id", 0, "The ID of the token that you own")
	cmd.MarkFlagRequired("token-id")

	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "(Optional) The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "(Optional) The HTTP JSON API server's ip-address")

	return cmd
}

func doBurnToken() {
	//
	// STEP 1:
	// Get our project dependencies in order.
	//
	logger := logger.NewProvider()

	//
	// STEP 2:
	// Create our request payload.
	//

	httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, burnTokensURL)

	metadata := handler.BurnTokenRequestIDO{
		TokenOwnerAddress:  flagBurnTokenOwnerAddress,
		TokenOwnerPassword: flagBurnTokenOwnerPassword,
		TokenID:            flagBurnTokenID,
	}
	logger.Debug("Burning token between accounts in blockchain",
		slog.Any("node-url", httpEndpoint),
		slog.Any("token-id", flagBurnTokenID),
	)

	//
	// STEP 3
	// Convert request to binary and submit to running HTTP JSON API.
	//

	reqBytes, err := json.Marshal(&metadata)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	if reqBytes == nil {
		log.Fatal("nothing marshalled")
	}

	//
	// STEP 4
	// Setup our HTTP client for sending
	//

	r, err := http.NewRequest("POST", httpEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("failed to setup post request: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "POST"))

	//
	// STEP 5
	// Wait for the submission to finish sending and then get the resposne.
	//

	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		//
		// STEP 6
		// If successful then return our results.
		//
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
	logger.Debug("Token burn request submitted successful")
}
