package http

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/rpc/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/service"
)

type RPCServer interface {
	Run()
	Shutdown()
}

type rpcServerImpl struct {
	cfg         *config.Config
	logger      *slog.Logger
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
}

func NewRPCServer(
	cfg *config.Config,
	logger *slog.Logger,
	getKeyService *service.GetKeyService,
) RPCServer {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.App.RPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new RPC server
	accountServer := handler.NewAccountServer(cfg, logger, getKeyService)

	// Register RPC server
	rpc.Register(accountServer)
	rpc.HandleHTTP()

	port := &rpcServerImpl{
		cfg:     cfg,
		logger:  logger,
		tcpAddr: tcpAddr,
	}
	return port
}

func (port *rpcServerImpl) Run() {
	// ctx := context.Background()
	port.logger.Info("Running TCP RPC server",
		slog.Any("listen_addr", port.tcpAddr))

	// Listen for requests on `address:port`.
	l, err := net.ListenTCP("tcp", port.tcpAddr)
	if err != nil {
		port.logger.Error("listen failed", slog.Any("error", err))

		// DEVELOPERS NOTE: We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of app.
		log.Fatalf("failed to listen and server: %v", err)
	}

	// Save our TCP listener so we can close later on graceful shutdown.
	port.tcpListener = l

	// Safety net for 'too many open files' issue on legacy code.
	// Set a sane timeout duration for the http.DefaultClient, to ensure idle connections are terminated.
	// Reference: https://stackoverflow.com/questions/37454236/net-http-server-too-many-open-files-error
	http.DefaultClient.Timeout = time.Minute * 10

	// DEVELOPER NOTES:
	// If you get "too many open files" then please read the following article
	// http://publib.boulder.ibm.com/httpserv/ihsdiag/too_many_open_files.html
	// so you can run in your console:
	// $ ulimit -H -n 4096
	// $ ulimit -n 4096

	http.Serve(l, nil)
}

func (port *rpcServerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down TCP RPC server")
	port.tcpListener.Close()
}
