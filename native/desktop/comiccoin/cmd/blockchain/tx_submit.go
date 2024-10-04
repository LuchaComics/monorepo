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
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
)

func submitTxCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "submit",
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network",
		Run: func(cmd *cobra.Command, args []string) {
			doSubmitTransactionCommand()
		},
	}

	cmd.Flags().StringVar(&flagAccountAddress, "sender-account-id", "", "The id of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("sender-account-id")

	cmd.Flags().StringVar(&flagPassword, "sender-account-password", "", "The password to unlock the account which will transfer the coin")
	cmd.MarkFlagRequired("sender-account-password")

	cmd.Flags().Uint64Var(&flagAmount, "value", 0, "The amount of coins to send")
	cmd.MarkFlagRequired("value")

	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The name of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("recipient-address")

	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}

func doSubmitTransactionCommand() {
	logger := logger.NewProvider()

	//
	// Create our request payload.
	//

	httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, transactionsURL)

	req := &httphandler.CreateTransactionRequestIDO{
		SenderAccountAddress:  flagAccountAddress,
		SenderAccountPassword: flagPassword,
		RecipientAddress:      flagRecipientAddress,
		Value:                 flagAmount,
		Data:                  nil,
	}

	logger.Debug("Submitting to blockchain",
		slog.Any("node-url", httpEndpoint),
		slog.Any("sender-account-addresss", flagAccountAddress),
		slog.Any("sender-account-password", flagPassword),
		slog.Any("value", flagAmount),
		slog.Any("recipient-address", flagRecipientAddress),
		slog.Any("request", req),
	)

	//
	// Convert request to binary and submit to running HTTP JSON API.
	//

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	if reqBytes == nil {
		log.Fatal("nothing marshalled")
	}
	r, err := http.NewRequest("POST", httpEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("failed to setup post request: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "POST"))

	//
	// Wait for the submission to finish sending and then get the resposne.
	//

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
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

	logger.Debug("Pending transaction submitted successful")
}
