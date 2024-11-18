package rpc

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"

	rpchandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/rpc/handler"
)

// RPCServer represents an RPC server that handles incoming requests.
type RPCServer interface {
	// Run starts the RPC server.
	Run(address, port string)

	// Shutdown shuts down the RPC server.
	Shutdown()
}

// RPCServerImpl is an implementation of the RPCServer interface.
type RPCServerImpl struct {
	// logger is the logger for the RPC server.
	logger *slog.Logger

	rpcApi *rpchandler.ComicCoinRPCServer
}

// NewRPCServer creates a new RPC server instance.
func NewRPCServer(
	logger *slog.Logger,
) RPCServer {
	// Create a new RPC server
	myServer := rpchandler.NewComicCoinRPCServer(logger)

	// Create a new RPC server instance.
	port := &RPCServerImpl{
		logger: logger,
		rpcApi: myServer,
	}

	return port
}

// Run starts the RPC server.
func (impl *RPCServerImpl) Run(address, port string) {
	// Log a message to indicate that the RPC server is running.
	impl.logger.Info("Running RPC API",
		slog.String("listen_address", port))

	// Register RPC server
	rpc.Register(impl.rpcApi)
	rpc.HandleHTTP()
	// Listen for requests on port
	l, e := net.Listen("tcp", fmt.Sprintf("%v:%v", address, port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}

// Shutdown shuts down the RPC server.
func (port *RPCServerImpl) Shutdown() {
	// Log a message to indicate that the RPC server is shutting down.
	port.logger.Info("Gracefully shutting down RPC API")
}
