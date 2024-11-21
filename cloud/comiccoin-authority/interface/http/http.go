package http

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/handler"
	mid "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/middleware"
)

// HTTPServer represents an HTTP server that handles incoming requests.
type HTTPServer interface {
	// Run starts the HTTP server.
	Run()

	// Shutdown shuts down the HTTP server.
	Shutdown()
}

// httpServerImpl is an implementation of the HTTPServer interface.
type httpServerImpl struct {
	// cfg is the configuration for the HTTP server.
	cfg *config.Configuration

	// logger is the logger for the HTTP server.
	logger *slog.Logger

	// server is the underlying HTTP server.
	server *http.Server

	getVersionHTTPHandler                                     *handler.GetVersionHTTPHandler
	getHealthCheckHTTPHandler                                 *handler.GetHealthCheckHTTPHandler
	getGenesisBlockDataHTTPHandler                            *handler.GetGenesisBlockDataHTTPHandler
	getBlockchainStateHTTPHandler                             *handler.GetBlockchainStateHTTPHandler
	getBlockDataHTTPHandler                                   *handler.GetBlockDataHTTPHandler
	signedTransactionSubmissionHTTPHandler                    *handler.SignedTransactionSubmissionHTTPHandler
	mempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler *handler.MempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler
	blockchainStateChangeEventDTOHTTPHandler                  *handler.BlockchainStateChangeEventDTOHTTPHandler
}

// NewHTTPServer creates a new HTTP server instance.
func NewHTTPServer(
	cfg *config.Configuration,
	logger *slog.Logger,
	mid mid.Middleware,
	http1 *handler.GetVersionHTTPHandler,
	http2 *handler.GetHealthCheckHTTPHandler,
	http3 *handler.GetGenesisBlockDataHTTPHandler,
	http4 *handler.GetBlockchainStateHTTPHandler,
	http5 *handler.BlockchainStateChangeEventDTOHTTPHandler,
	http6 *handler.GetBlockDataHTTPHandler,
	http7 *handler.SignedTransactionSubmissionHTTPHandler,
	http8 *handler.MempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler,
) HTTPServer {
	// Check if the HTTP address is set in the configuration.
	if cfg.App.IP == "" {
		log.Fatal("http server missing ip address")
	}
	if cfg.App.Port == "" {
		log.Fatal("http server missing port number")
	}

	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// Set up CORS middleware to allow all origins.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.App.IP, cfg.App.Port),
		Handler: handler,
	}

	// Create a new HTTP server instance.
	port := &httpServerImpl{
		cfg:                                      cfg,
		logger:                                   logger,
		server:                                   srv,
		getVersionHTTPHandler:                    http1,
		getHealthCheckHTTPHandler:                http2,
		getGenesisBlockDataHTTPHandler:           http3,
		getBlockchainStateHTTPHandler:            http4,
		blockchainStateChangeEventDTOHTTPHandler: http5,
		getBlockDataHTTPHandler:                  http6,
		signedTransactionSubmissionHTTPHandler:   http7,
		mempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler: http8,
	}

	// Attach the HTTP server controller to the ServeMux.
	mux.HandleFunc("/", mid.Attach(port.HandleRequests))

	return port
}

// Run starts the HTTP server.
func (port *httpServerImpl) Run() {
	// Log a message to indicate that the HTTP server is running.
	port.logger.Info("Running HTTP JSON API",
		slog.String("listen_address", fmt.Sprintf("%v:%v", port.cfg.App.IP, port.cfg.App.Port)))

	// Start the HTTP server.
	if err := port.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Log an error message if the server fails to start.
		port.logger.Error("listen failed", slog.Any("error", err))

		// Terminate the application if the server fails to start.
		log.Fatalf("failed to listen and server: %v", err)
	}
}

// Shutdown shuts down the HTTP server.
func (port *httpServerImpl) Shutdown() {
	// Log a message to indicate that the HTTP server is shutting down.
	port.logger.Info("Gracefully shutting down HTTP JSON API")
}

// HandleRequests handles incoming HTTP requests.
func (port *httpServerImpl) HandleRequests(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to application/json.
	w.Header().Set("Content-Type", "application/json")

	// Get the URL path tokens from the request context.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)

	// Log a message to indicate that a request has been received.
	// But only do this if client is attempting to access our API endpoints.
	if n > 2 {
		port.logger.Debug("",
			slog.Any("method", r.Method),
			slog.Any("url_tokens", p),
			slog.Int("url_token_count", n))
	}

	// Handle the request based on the URL path tokens.
	switch {
	case n == 1 && p[0] == "version" && r.Method == http.MethodGet:
		port.getVersionHTTPHandler.Execute(w, r)
	case n == 1 && p[0] == "health-check" && r.Method == http.MethodGet:
		port.getHealthCheckHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "genesis" && r.Method == http.MethodGet:
		port.getGenesisBlockDataHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "blockchain-state" && r.Method == http.MethodGet:
		port.getBlockchainStateHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "api" && p[1] == "v1" && p[2] == "blockchain-state" && p[3] == "changes" && r.Method == http.MethodGet:
		port.blockchainStateChangeEventDTOHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "api" && p[1] == "v1" && p[2] == "blockdata" && r.Method == http.MethodGet:
		port.getBlockDataHTTPHandler.Execute(w, r, p[3])
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "signed-transaction" && r.Method == http.MethodPost:
		port.signedTransactionSubmissionHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "mempool-transactions" && r.Method == http.MethodPost:
		port.mempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler.Execute(w, r)
	// --- CATCH ALL: D.N.E. ---
	default:
		// Log a message to indicate that the request is not found.
		port.logger.Debug("404 request",
			slog.Any("method", r.Method),
			slog.Any("url_tokens", p),
			slog.Int("url_token_count", n),
		)

		// Return a 404 response.
		http.NotFound(w, r)
	}
}
