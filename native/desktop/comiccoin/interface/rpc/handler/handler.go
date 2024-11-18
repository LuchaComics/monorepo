package handler

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type ComicCoinRPCServer struct {
	logger                               *slog.Logger
	getOrDownloadNonFungibleTokenService *service.GetOrDownloadNonFungibleTokenService
}

func NewComicCoinRPCServer(
	logger *slog.Logger,
	s1 *service.GetOrDownloadNonFungibleTokenService,
) *ComicCoinRPCServer {

	// Create a new RPC server instance.
	port := &ComicCoinRPCServer{
		logger:                               logger,
		getOrDownloadNonFungibleTokenService: s1,
	}

	return port
}
