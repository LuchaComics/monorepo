package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	a_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

func init() {
	rootCmd.AddCommand(httpJsonApiCmd)
	httpJsonApiCmd.AddCommand(httpJsonApiNewAccountCmd())
}

var httpJsonApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Execute commands for local running ComicCoin node instance via HTTP JSON API",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

// HTTP endpoints
const (
	accountsURL = "/v1/api/accounts"
)

func httpJsonApiNewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new-account",
		Short: "Creates a new wallet in our ComicCoin node local filesystem and encrypts it with the inputted password",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()
			logger.Debug("Creating new account...")

			httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, accountsURL)

			req := &a_c.AccountCreateRequestIDO{
				Name:           "test",
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

			post := &a_c.AccountDetailResponseIDO{}
			derr := json.NewDecoder(res.Body).Decode(post)
			if derr != nil {
				log.Fatalf("failed to decode response: %v", err)
			}

			if res.StatusCode != http.StatusCreated {
				log.Fatal("failed to get created status")
			}

			logger.Debug("Account created",
				slog.String("name", post.Name),
				slog.String("filepath", post.WalletFilepath),
				slog.String("address", post.WalletAddress),
			)
		},
	}

	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to encrypt the new wallet with")
	cmd.MarkFlagRequired("password")
	cmd.Flags().IntVar(&flagListenHTTPPort, "http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}
