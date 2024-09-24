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

	a_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

var (
	flagWalletAddress string
)

// HTTP endpoints
const (
	ledgerBalanceURL = "/v1/api/ledger/${ACCOUNT_NAME}/balance"
)

func httpJsonApiLedgerGetBalanceCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "balance",
		Short: "Get balance of the account in the ledger",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("Get ledger detail...")

			modifiedBalanceDetailURL := strings.ReplaceAll(ledgerBalanceURL, "${ACCOUNT_NAME}", flagAccountName)
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

			post := &a_c.LedgerBalanceResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Ledger balance retrieved",
				slog.Any("amount", post.Amount),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountName, "account-name", "", "The name of the account we want to lookup in the ledger to get our balance for")
	cmd.MarkFlagRequired("account-name")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
