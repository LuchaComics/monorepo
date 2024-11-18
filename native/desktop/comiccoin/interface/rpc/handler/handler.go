package handler

import "log/slog"

type ComicCoinRPCServer struct {
	// logger is the logger for the RPC server.
	logger *slog.Logger
}

func NewComicCoinRPCServer(
	logger *slog.Logger,
) *ComicCoinRPCServer {

	// Create a new RPC server instance.
	port := &ComicCoinRPCServer{
		logger: logger,
	}

	return port
}
