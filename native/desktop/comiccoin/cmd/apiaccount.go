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

	a_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

// HTTP endpoints
const (
	accountsURL      = "/v1/api/accounts"
	accountDetailURL = "/v1/api/account/${ACCOUNT_NAME}"
)

func httpJsonApiNewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new wallet in our ComicCoin node local filesystem and encrypts it with the inputted password",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("Creating new account...")

			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, accountsURL)

			req := &a_c.AccountCreateRequestIDO{
				Name:           flagAccountName,
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

			post := &a_c.AccountDetailResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Account created",
				slog.String("name", post.Name),
				slog.String("address", post.WalletAddress),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountName, "account-name", "", "The name to assign this account")
	cmd.MarkFlagRequired("account-name")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to encrypt the new wallet with")
	cmd.MarkFlagRequired("password")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}

func httpJsonApiGetAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get account detail",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("Get account detail...")

			modifiedBalanceDetailURL := strings.ReplaceAll(accountDetailURL, "${ACCOUNT_NAME}", flagAccountName)
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

			post := &a_c.AccountDetailResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("Account retrieved",
				slog.String("name", post.Name),
				slog.String("address", post.WalletAddress),
			)
		},
	}

	cmd.Flags().StringVar(&flagAccountName, "name", "", "The name to lookup this account")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}

func httpJsonApiListAccountsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all local accounts",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("List accounts...")

			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, accountsURL)

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

			posts := &[]a_c.AccountDetailResponseIDO{}
			if err := json.NewDecoder(teeReader).Decode(&posts); err != nil {
				logger.Error("decoding string error",
					slog.Any("err", err),
					slog.String("json", rawJSON.String()),
				)
				return
			}

			logger.Debug("",
				slog.Any("posts", posts),
			)
		},
	}

	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
