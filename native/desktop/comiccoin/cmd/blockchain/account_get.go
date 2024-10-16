package blockchain

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

	ah "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/interface/http/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

func httpJsonApiGetAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get account detail",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			logger.Debug("Creating get account...")

			logger.Debug("Get account detail...")

			modifiedAccountDetailURL := strings.ReplaceAll(accountDetailURL, "${ACCOUNT_ADDRESS}", flagAccountAddress)
			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, modifiedAccountDetailURL)

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

			post := &ah.AccountGetResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Account retrieved",
				slog.Any("balance", post.Balance),
				slog.Any("nonce", post.Nonce),
				slog.String("address", post.Address),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountAddress, "account-address", "", "The address value to lookup the account by")
	cmd.MarkFlagRequired("account-address")
	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
