package tokens

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

// Command line argument flags
var (
	flagAccountAddress string
)

func ListTokensCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List the tokens owned by an account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunListTokensCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagAccountAddress, "owner-address", "", "The address value to lookup the tokens by")
	cmd.MarkFlagRequired("owner-address")

	return cmd
}

func doRunListTokensCommand() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()
	accountAddress := common.HexToAddress(strings.ToLower(flagAccountAddress))

	logger.Debug("Retrieving tokens...",
		slog.Any("addresss", accountAddress))

	toks, retrieveErr := rpcClient.ListTokensByOwnerAddress(
		ctx,
		&accountAddress,
	)
	if retrieveErr != nil {
		log.Fatalf("Failed execute get token: %v", retrieveErr)
	}
	if toks == nil {
		logger.Warn("Tokens does not exist",
			slog.Any("address", flagAccountAddress))
		return
	}
	for _, tok := range toks {
		logger.Debug("Token",
			slog.Any("ID", tok.GetID()),
			slog.Any("MetadataURI", tok.MetadataURI),
		)
	}
}
