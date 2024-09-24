package http

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	account_http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/httptransport"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
	mid "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/http/middleware"
)

type httpInputPort struct {
	cfg     *config.Config
	logger  *slog.Logger
	server  *http.Server
	account *account_http.Handler
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	mid mid.Middleware,
	acc *account_http.Handler,
) inputport.InputPortServer {
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
	port := &httpInputPort{
		cfg:     cfg,
		logger:  logger,
		server:  srv,
		account: acc,
	}

	// Attach the HTTP server controller to the ServerMux.
	mux.HandleFunc("/", mid.Attach(port.HandleRequests))

	return port
}

func (port *httpInputPort) Run() {
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

func (port *httpInputPort) Shutdown() {
	port.logger.Info("Gracefully shutting down HTTP JSON API")
}

func (port *httpInputPort) HandleRequests(w http.ResponseWriter, r *http.Request) {
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
	case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodGet:
		port.account.List(w, r)
	case n == 3 && p[0] == "v1" && p[1] == "api" && p[2] == "accounts" && r.Method == http.MethodPost:
		port.account.Create(w, r)
	case n == 4 && p[0] == "v1" && p[1] == "api" && p[2] == "account" && r.Method == http.MethodGet:
		port.account.GetByName(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodPut:
	// 	port.ComicSubmission.UpdateByID(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodDelete:
	// 	port.ComicSubmission.ArchiveByID(w, r, p[3])
	// case n == 5 && p[1] == "v1" && p[2] == "comic-submission" && p[4] == "perma-delete" && r.Method == http.MethodDelete:
	// 	port.ComicSubmission.DeleteByID(w, r, p[3])
	// --- CATCH ALL: D.N.E. ---
	default:
		port.logger.Debug("404 request",
			slog.Any("t[0]", p[0]),
			slog.Any("t[1]", p[1]),
			slog.Any("t[2]", p[2]),
			slog.Any("method", r.Method),
			slog.Any("url_tokens", p),
			slog.Int("url_token_count", n),
		)
		http.NotFound(w, r)
	}
}
