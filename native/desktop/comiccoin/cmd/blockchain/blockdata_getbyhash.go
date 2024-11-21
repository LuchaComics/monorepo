package blockchain

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

var (
	flagHash string
)

func BlockDataGetByHashCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockdata",
		Short: "Get account details",
		Run: func(cmd *cobra.Command, args []string) {
			doRunBlockDataGetByHash()
		},
	}

	cmd.Flags().StringVar(&flagHash, "hash", "", "The hash value to lookup the blockdata by")
	cmd.MarkFlagRequired("hash")

	return cmd
}

func doRunBlockDataGetByHash() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	blockData, err := rpcClient.GetBlockDataByHash(ctx, flagHash)
	if err != nil {
		log.Fatalf("Failed to get block data: %v\n", err)
	}

	if blockData == nil {
		err := errors.New("Block data does not exist")
		log.Fatalf("Failed to getting block data: %v\n", err)
	}

	logger.Info("Block data retrieved",
		slog.Any("Hash", blockData.Hash),
	)
}
