package http

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http/handler"
	mid "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http/middleware"
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

	// transferCoinHTTPHandler is the handler for creating transactions.
	transferCoinHTTPHandler *handler.TransferCoinHTTPHandler

	// mintTokenHTTPHandler is the handler for minting new Token.
	mintTokenHTTPHandler *handler.ProofOfAuthorityTokenMintHTTPHandler

	// transferTokenHTTPHandler is the handler for transfering Tokens between accounts.
	transferTokenHTTPHandler *handler.TransferTokenHTTPHandler

	// burnTokenHTTPHandler is the handler for burning Tokens.
	burnTokenHTTPHandler *handler.BurnTokenHTTPHandler

	// getTokenHTTPHandler is the handler for getting Token detail.
	getTokenHTTPHandler *handler.GetTokenHTTPHandler
}

// NewHTTPServer creates a new HTTP server instance.
func NewHTTPServer(
	cfg *config.Config,
	logger *slog.Logger,
	mid mid.Middleware,
	createAccountHTTPHandler *handler.CreateAccountHTTPHandler,
	getAccountHTTPHandler *handler.GetAccountHTTPHandler,
	transferCoinHTTPHandler *handler.TransferCoinHTTPHandler,
	mintTokenHTTPHandler *handler.ProofOfAuthorityTokenMintHTTPHandler,
	transferTokenHTTPHandler *handler.TransferTokenHTTPHandler,
	burnTokenHTTPHandler *handler.BurnTokenHTTPHandler,
	getTokenHTTPHandler *handler.GetTokenHTTPHandler,
) HTTPServer {
	// Check if the HTTP address is set in the configuration.
	if cfg.App.HTTPAddress == "" {
		log.Fatal("NewHTTPServer: missing http address")
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
		cfg:                      cfg,
		logger:                   logger,
		server:                   srv,
		createAccountHTTPHandler: createAccountHTTPHandler,
		getAccountHTTPHandler:    getAccountHTTPHandler,
		transferCoinHTTPHandler:  transferCoinHTTPHandler,
		mintTokenHTTPHandler:     mintTokenHTTPHandler,
		transferTokenHTTPHandler: transferTokenHTTPHandler,
		burnTokenHTTPHandler:     burnTokenHTTPHandler,
		getTokenHTTPHandler:      getTokenHTTPHandler,
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
		log.Fatalf("httpServerImpl: Run: failed to listen and server: %v", err)
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

	// --- COINS --- //
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "coins" && p[3] == "transfer" && r.Method == http.MethodPost:
		// Handle the request to create a transaction.
		port.transferCoinHTTPHandler.Execute(w, r)

		// --- TOKENS --- //
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "token" && r.Method == http.MethodGet:
		// Handle the request to getting a transaction.
		port.getTokenHTTPHandler.Execute(w, r, p[3])
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "tokens" && p[3] == "mint" && r.Method == http.MethodPost:
		// Handle the request to create a transaction.
		port.mintTokenHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "tokens" && p[3] == "transfer" && r.Method == http.MethodPost:
		// Handle the request to create a transaction.
		port.transferTokenHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "tokens" && p[3] == "burn" && r.Method == http.MethodPost:
		// Handle the request to create a transaction.
		port.burnTokenHTTPHandler.Execute(w, r)

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