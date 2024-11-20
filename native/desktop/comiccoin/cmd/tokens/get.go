package tokens

import (
	"context"
	"log"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

// Command line argument flags
var (
	flagTokenID string
)

func GetTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get the token",
		Run: func(cmd *cobra.Command, args []string) {
			doRunGetTokenCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The unique token identification to use to lookup the token")
	cmd.MarkFlagRequired("token-id")

	return cmd
}

func doRunGetTokenCommand() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()
	tokenID, ok := new(big.Int).SetString(flagTokenID, 10)
	if !ok {
		log.Fatal("Failed convert `token_id` to big.Int")
	}

	logger.Debug("Token retrieving...",
		slog.Any("token_id", tokenID))

	tok, retrieveErr := rpcClient.GetToken(
		ctx,
		tokenID,
	)
	if retrieveErr != nil {
		log.Fatalf("Failed execute get token: %v", retrieveErr)
	}
	if tok == nil {
		logger.Warn("Token does not exist",
			slog.Any("token_id", tokenID))
	} else {
		logger.Debug("Retrieved token",
			slog.Any("token_id", tokenID),
			slog.Any("owner", tok.Owner),
			slog.Any("metadata_uri", tok.MetadataURI),
			slog.Any("nonce", tok.GetNonce()))
	}
}
