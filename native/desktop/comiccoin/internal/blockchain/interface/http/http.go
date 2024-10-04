package http

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
	mid "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/middleware"
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
	cfg *config.Config

	// logger is the logger for the HTTP server.
	logger *slog.Logger

	// server is the underlying HTTP server.
	server *http.Server

	// createAccountHTTPHandler is the handler for creating accounts.
	createAccountHTTPHandler *handler.CreateAccountHTTPHandler

	// getAccountHTTPHandler is the handler for getting accounts.
	getAccountHTTPHandler *handler.GetAccountHTTPHandler

	// createTransactionHTTPHandler is the handler for creating transactions.
	createTransactionHTTPHandler *handler.CreateTransactionHTTPHandler
}

// NewHTTPServer creates a new HTTP server instance.
func NewHTTPServer(
	// cfg is the configuration for the HTTP server.
	cfg *config.Config,

	// logger is the logger for the HTTP server.
	logger *slog.Logger,

	// mid is the middleware for the HTTP server.
	mid mid.Middleware,

	// createAccountHTTPHandler is the handler for creating accounts.
	createAccountHTTPHandler *handler.CreateAccountHTTPHandler,

	// getAccountHTTPHandler is the handler for getting accounts.
	getAccountHTTPHandler *handler.GetAccountHTTPHandler,

	// createTransactionHTTPHandler is the handler for creating transactions.
	createTransactionHTTPHandler *handler.CreateTransactionHTTPHandler,
) HTTPServer {
	// Check if the HTTP address is set in the configuration.
	if cfg.App.HTTPAddress == "" {
		log.Fatal("missing http address")
	}

	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// Set up CORS middleware to allow all origins.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	srv := &http.Server{
		Addr:    cfg.App.HTTPAddress,
		Handler: handler,
	}

	// Create a new HTTP server instance.
	port := &httpServerImpl{
		cfg:                          cfg,
		logger:                       logger,
		server:                       srv,
		createAccountHTTPHandler:     createAccountHTTPHandler,
		getAccountHTTPHandler:        getAccountHTTPHandler,
		createTransactionHTTPHandler: createTransactionHTTPHandler,
	}

	// Attach the HTTP server controller to the ServeMux.
	mux.HandleFunc("/", mid.Attach(port.HandleRequests))

	return port
}

// Run starts the HTTP server.
func (port *httpServerImpl) Run() {
	// Log a message to indicate that the HTTP server is running.
	port.logger.Info("Running HTTP JSON API",
		slog.String("listen_address", port.cfg.App.HTTPAddress))

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
	port.logger.Info("",
		slog.Any("method", r.Method),
		slog.Any("url_tokens", p),
		slog.Int("url_token_count", n))

	// Handle the request based on the URL path tokens.
	switch {
	// --- ACCOUNTS --- //
	// case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodGet:
	// 	port.account.List(w, r)
	case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodPost:
		// Handle the request to create an account.
		port.createAccountHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "account" && r.Method == http.MethodGet:
		// Handle the request to get an account.
		port.getAccountHTTPHandler.Execute(w, r, p[3])
		// // case n == 4 && p[0] == "v1" && p[1] == "account" && r.Method == http.MethodPut:
		// // 	port.account.UpdateByName(w, r, p[3])
		// case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "account" && r.Method == http.MethodDelete:
		// 	port.account.DeleteByName(w, r, p[3])

		// --- TRANSACTIONS --- //
	case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "txs" && r.Method == http.MethodPost:
		// Handle the request to create a transaction.
		port.createTransactionHTTPHandler.Execute(w, r)

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
