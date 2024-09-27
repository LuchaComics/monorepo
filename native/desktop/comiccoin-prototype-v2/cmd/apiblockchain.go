package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	a_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

// HTTP endpoints
const (
	blockchainBalanceURL = "/v1/api/blockchain/${ACCOUNT_NAME}/balance"
	blockchainSubmitURL  = "/v1/api/blockchain/submit"
)

func httpJsonApiBlockchainGetBalanceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "balance",
		Short: "Get balance of the account in the blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("Get blockchain detail...")

			modifiedBalanceDetailURL := strings.ReplaceAll(blockchainBalanceURL, "${ACCOUNT_NAME}", flagAccountName)
			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, modifiedBalanceDetailURL)

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

			var rawJSON bytes.Buffer
			teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

			post := &a_c.BlockchainBalanceResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Blockchain balance retrieved",
				slog.Any("amount", post.Amount),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountName, "account-name", "", "The name of the account we want to lookup in the blockchain to get our balance for")
	cmd.MarkFlagRequired("account-name")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}

func httpJsonApiBlockchainSubmitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "submit",
		Short: "Submit a pending transaction to the ComicCoin blockchain network to mine, verify and then approve globally in the network",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()

			//
			// Create our request payload.
			//

			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, blockchainSubmitURL)

			req := &a_c.BlockchainSubmitRequestIDO{
				FromAccountName:       flagAccountName,
				AccountWalletPassword: flagPassword,
				To:                    flagRecipientAddress,
				Value:                 flagAmount,
				Data:                  nil,
			}

			logger.Debug("Submitting to blockchain",
				slog.Any("node-url", httpEndpoint),
				slog.Any("sender-account-name", flagAccountName),
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
		},
	}

	cmd.Flags().StringVar(&flagAccountName, "sender-account-name", "", "The name of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("sender-account-name")
	cmd.Flags().StringVar(&flagPassword, "sender-account-password", "", "The password to unlock the account which will transfer the coin")
	cmd.MarkFlagRequired("sender-account-password")
	cmd.Flags().Uint64Var(&flagAmount, "value", 0, "The amount of coins to send")
	cmd.MarkFlagRequired("value")
	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The name of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("recipient-address")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
