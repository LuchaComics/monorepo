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

	ah "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

func httpJsonApiNewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new wallet in our ComicCoin node local filesystem and encrypts it with the inputted password",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			logger.Debug("Creating new account...")

			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, accountsURL)

			req := &ah.AccountCreateRequestIDO{
				ID:             flagAccountID,
				WalletPassword: flagPassword,
			}

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

			client := &http.Client{}
			res, err := client.Do(r)
			if err != nil {
				log.Fatalf("failed to do post request: %v", err)
			}

			defer res.Body.Close()

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

			post := &ah.AccountCreateResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Account created",
				slog.String("id", post.ID),
				slog.String("address", post.WalletAddress),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountID, "account-id", "", "The ID to assign this account")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVar(&flagPassword, "wallet-password", "", "The password to encrypt the new wallet with")
	cmd.MarkFlagRequired("wallet-password")
	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
