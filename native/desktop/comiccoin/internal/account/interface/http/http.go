package http

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/http/handler"
	mid "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/http/middleware"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type HTTPServer interface {
	Run()
	Shutdown()
}

type httpServerImpl struct {
	cfg                      *config.Config
	logger                   *slog.Logger
	server                   *http.Server
	createAccountHTTPHandler *handler.CreateAccountHTTPHandler
}

func NewHTTPServer(
	cfg *config.Config,
	logger *slog.Logger,
	mid mid.Middleware,
	createAccountHTTPHandler *handler.CreateAccountHTTPHandler,
) HTTPServer {
	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	addr := fmt.Sprintf("%s:%d", cfg.App.HTTPIP, cfg.App.HTTPPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	port := &httpServerImpl{
		cfg:                      cfg,
		logger:                   logger,
		server:                   srv,
		createAccountHTTPHandler: createAccountHTTPHandler,
	}

	// Attach the HTTP server controller to the ServerMux.
	mux.HandleFunc("/", mid.Attach(port.HandleRequests))

	return port
}

func (port *httpServerImpl) Run() {
	// ctx := context.Background()
	port.logger.Info("Running HTTP JSON API",
		slog.Int("listen_port", port.cfg.App.HTTPPort),
		slog.String("listen_ip", port.cfg.App.HTTPIP))
	if err := port.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		port.logger.Error("listen failed", slog.Any("error", err))

		// DEVELOPERS NOTE: We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of app.
		log.Fatalf("failed to listen and server: %v", err)
	}
}

func (port *httpServerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down HTTP JSON API")
}

func (port *httpServerImpl) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)

	port.logger.Info("",
		slog.Any("method", r.Method),
		slog.Any("url_tokens", p),
		slog.Int("url_token_count", n))

	switch {
	// // --- ACCOUNTS --- //
	// case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodGet:
	// 	port.account.List(w, r)
	case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodPost:
		port.createAccountHTTPHandler.Execute(w, r)
	// case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "account" && r.Method == http.MethodGet:
	// 	port.account.GetByName(w, r, p[3])
	// // case n == 4 && p[0] == "v1" && p[1] == "account" && r.Method == http.MethodPut:
	// // 	port.account.UpdateByName(w, r, p[3])
	// case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "account" && r.Method == http.MethodDelete:
	// 	port.account.DeleteByName(w, r, p[3])

	// --- CATCH ALL: D.N.E. ---
	default:
		port.logger.Debug("404 request",
			slog.Any("method", r.Method),
			slog.Any("url_tokens", p),
			slog.Int("url_token_count", n),
		)
		http.NotFound(w, r)
	}
}
