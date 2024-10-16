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
	handler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http/handler"
)

// Command line argument flags
var (
	flagProofOfAuthorityAccountAddress string
	flagProofOfAuthorityWalletPassword string
	flagMintRecipientAddress           string
	flagMintMetadataURI                string
)

func httpJsonApiMintTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mint",
		Short: "(PoA & Authority only) Creates a new non-fungible token in our blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			doMintToken()
		},
	}

	cmd.Flags().StringVar(&flagProofOfAuthorityAccountAddress, "poa-address", "", "(Required for `PoA` consensus protocol) The address of the authority's account")
	cmd.MarkFlagRequired("poa-address")
	cmd.Flags().StringVar(&flagProofOfAuthorityWalletPassword, "poa-password", "", "(Required for `PoA` consensus protocol) The password in the authority's wallet")
	cmd.MarkFlagRequired("poa-password")
	cmd.Flags().StringVar(&flagMintRecipientAddress, "recipient-address", "", "The address of the account whom will receive this Token")
	cmd.MarkFlagRequired("recipient-address")

	// Fields for inputting the Token
	cmd.Flags().StringVar(&flagMintMetadataURI, "metadata-uri", "", "The location of this tokens metadata file.")
	cmd.MarkFlagRequired("metadata-uri")

	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "(Optional) The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "(Optional) The HTTP JSON API server's ip-address")

	return cmd
}

func doMintToken() {
	//
	// STEP 1:
	// Get our project dependencies in order.
	//
	logger := logger.NewProvider()

	//
	// STEP 2:
	// Create our request payload.
	//

	httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, mintTokensURL)

	metadata := handler.ProofOfAuthorityTokenMintRequestIDO{
		ProofOfAuthorityAccountAddress: flagProofOfAuthorityAccountAddress,
		ProofOfAuthorityWalletPassword: flagProofOfAuthorityWalletPassword,
		To:                             flagMintRecipientAddress,
		MetadataURI:                    flagMintMetadataURI,
	}
	logger.Debug("Creating new Token in blockchain",
		slog.Any("node-url", httpEndpoint),
		slog.Any("metadata_uri", flagMintMetadataURI),
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
	logger.Debug("Token mint request submitted successful")
}
