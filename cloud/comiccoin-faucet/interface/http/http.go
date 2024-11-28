package http

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/handler"
	mid "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/middleware"
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

	getVersionHTTPHandler *handler.GetVersionHTTPHandler

	getHealthCheckHTTPHandler *handler.GetHealthCheckHTTPHandler

	gatewayRegisterCustomerHTTPHandler *handler.GatewayRegisterCustomerHTTPHandler
	gatewayLoginHTTPHandler            *handler.GatewayLoginHTTPHandler
	gatewayLogoutHTTPHandler           *handler.GatewayLogoutHTTPHandler
	gatewayRefreshTokenHTTPHandler     *handler.GatewayRefreshTokenHTTPHandler
	gatewayProfileDetailHTTPHandler    *handler.GatewayProfileDetailHTTPHandler
	gatewayProfileUpdateHTTPHandler    *handler.GatewayProfileUpdateHTTPHandler
	gatewayVerifyHTTPHandler           *handler.GatewayVerifyHTTPHandler
	gatewayChangePasswordHTTPHandler   *handler.GatewayChangePasswordHTTPHandler
	gatewayForgotPasswordHTTPHandler   *handler.GatewayForgotPasswordHTTPHandler
	gatewayResetPasswordHTTPHandler    *handler.GatewayResetPasswordHTTPHandler
}

// NewHTTPServer creates a new HTTP server instance.
func NewHTTPServer(
	cfg *config.Configuration,
	logger *slog.Logger,
	mid mid.Middleware,
	h1 *handler.GetVersionHTTPHandler,
	h2 *handler.GetHealthCheckHTTPHandler,
	h3 *handler.GatewayRegisterCustomerHTTPHandler,
	h4 *handler.GatewayLoginHTTPHandler,
	h5 *handler.GatewayLogoutHTTPHandler,
	h6 *handler.GatewayRefreshTokenHTTPHandler,
	h7 *handler.GatewayProfileDetailHTTPHandler,
	h8 *handler.GatewayProfileUpdateHTTPHandler,
	h9 *handler.GatewayVerifyHTTPHandler,
	h10 *handler.GatewayChangePasswordHTTPHandler,
	h11 *handler.GatewayForgotPasswordHTTPHandler,
	h12 *handler.GatewayResetPasswordHTTPHandler,
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
		cfg:                                cfg,
		logger:                             logger,
		server:                             srv,
		getVersionHTTPHandler:              h1,
		getHealthCheckHTTPHandler:          h2,
		gatewayRegisterCustomerHTTPHandler: h3,
		gatewayLoginHTTPHandler:            h4,
		gatewayLogoutHTTPHandler:           h5,
		gatewayRefreshTokenHTTPHandler:     h6,
		gatewayProfileDetailHTTPHandler:    h7,
		gatewayProfileUpdateHTTPHandler:    h8,
		gatewayVerifyHTTPHandler:           h9,
		gatewayChangePasswordHTTPHandler:   h10,
		gatewayForgotPasswordHTTPHandler:   h11,
		gatewayResetPasswordHTTPHandler:    h12,
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
	case n == 4 && p[0] == "api" && p[1] == "v1" && p[2] == "register" && p[3] == "customer" && r.Method == http.MethodPost:
		port.gatewayRegisterCustomerHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "login" && r.Method == http.MethodPost:
		port.gatewayLoginHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "logout" && r.Method == http.MethodPost:
		port.gatewayLogoutHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "refresh-token" && r.Method == http.MethodPost:
		port.gatewayRefreshTokenHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodGet:
		port.gatewayProfileDetailHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodPut:
		port.gatewayProfileUpdateHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "verify" && r.Method == http.MethodPost:
		port.gatewayVerifyHTTPHandler.Execute(w, r)
	case n == 4 && p[0] == "api" && p[1] == "v1" && p[2] == "profile" && p[3] == "change-password" && r.Method == http.MethodPut:
		port.gatewayChangePasswordHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "forgot-password" && r.Method == http.MethodPost:
		port.gatewayForgotPasswordHTTPHandler.Execute(w, r)
	case n == 3 && p[0] == "api" && p[1] == "v1" && p[2] == "password-reset" && r.Method == http.MethodPost:
		port.gatewayResetPasswordHTTPHandler.Execute(w, r)
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

/*
case n == 4 && p[1] == "v1" && p[2] == "register" && p[3] == "business" && r.Method == http.MethodPost:
	port.Gateway.RegisterBusiness(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "generate" && r.Method == http.MethodPost:
	port.Gateway.GenerateOTP(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "generate-qr-code" && r.Method == http.MethodPost:
	port.Gateway.GenerateOTPAndQRCodePNGImage(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "verify" && r.Method == http.MethodPost:
	port.Gateway.VerifyOTP(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "validate" && r.Method == http.MethodPost:
	port.Gateway.ValidateOTP(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "disable" && r.Method == http.MethodPost:
	port.Gateway.DisableOTP(w, r)
case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "recovery" && r.Method == http.MethodPost:
	port.Gateway.RecoveryOTP(w, r)

// --- REGISTRY --- // (TODO)
case n == 4 && p[1] == "v1" && p[2] == "cpsrn" && r.Method == http.MethodGet:
	port.ComicSubmission.GetRegistryByCPSRN(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "cpsrn" && p[4] == "qr-code" && r.Method == http.MethodGet:
	port.ComicSubmission.GetQRCodePNGImageOfRegisteryURLByCPSRN(w, r, p[3])

// --- SUBMISSIONS --- //
case n == 3 && p[1] == "v1" && p[2] == "comic-submissions" && r.Method == http.MethodGet:
	port.ComicSubmission.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "comic-submissions" && r.Method == http.MethodPost:
	port.ComicSubmission.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodGet:
	port.ComicSubmission.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodPut:
	port.ComicSubmission.UpdateByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodDelete:
	port.ComicSubmission.ArchiveByID(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "comic-submission" && p[4] == "perma-delete" && r.Method == http.MethodDelete:
	port.ComicSubmission.DeleteByID(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "operation" && p[4] == "set-customer" && r.Method == http.MethodPost:
	port.ComicSubmission.OperationSetCustomer(w, r)
case n == 5 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
	port.ComicSubmission.OperationCreateComment(w, r)
case n == 4 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "select-options" && r.Method == http.MethodGet:
	port.ComicSubmission.ListAsSelectOptionByFilter(w, r)
// case n == 5 && p[1] == "v1" && p[2] == "comic-submission" && p[4] == "file-attachments" && r.Method == http.MethodPost:
// 	port.ComicSubmission.CreateFileAttachment(w, r, p[3])

// --- ORGANIZATION --- //
case n == 3 && p[1] == "v1" && p[2] == "stores" && r.Method == http.MethodGet:
	port.Store.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "stores" && r.Method == http.MethodPost:
	port.Store.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "store" && r.Method == http.MethodGet:
	port.Store.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "store" && r.Method == http.MethodPut:
	port.Store.UpdateByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "store" && r.Method == http.MethodDelete:
	port.Store.DeleteByID(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "stores" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
	port.Store.OperationCreateComment(w, r)
case n == 4 && p[1] == "v1" && p[2] == "stores" && p[3] == "select-options" && r.Method == http.MethodGet:
	port.Store.ListAsSelectOptionByFilter(w, r)
case n == 4 && p[1] == "v1" && p[2] == "public" && p[3] == "stores-select-options" && r.Method == http.MethodGet:
	port.Store.PublicListAsSelectOptionByFilter(w, r)

// --- CUSTOMERS --- //
case n == 3 && p[1] == "v1" && p[2] == "customers" && r.Method == http.MethodGet:
	port.Customer.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "customers" && r.Method == http.MethodPost:
	port.Customer.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodGet:
	port.Customer.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodPut:
	port.Customer.UpdateByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodDelete:
	port.Customer.DeleteByID(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "customers" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
	port.Customer.OperationCreateComment(w, r)
case n == 5 && p[1] == "v1" && p[2] == "customers" && p[3] == "operation" && p[4] == "star" && r.Method == http.MethodPost:
	port.Customer.OperationStar(w, r)

// --- USERS --- //
case n == 3 && p[1] == "v1" && p[2] == "users" && r.Method == http.MethodGet:
	port.User.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "users" && r.Method == http.MethodPost:
	port.User.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodGet:
	port.User.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodPut:
	port.User.UpdateByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodDelete:
	port.User.DeleteByID(w, r, p[3])
case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
	port.User.OperationCreateComment(w, r)
case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operation" && p[4] == "star" && r.Method == http.MethodPost:
	port.User.OperationStar(w, r)
case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operation" && p[4] == "archive" && r.Method == http.MethodPost:
	port.User.OperationArchive(w, r)
case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operations" && p[4] == "change-password" && r.Method == http.MethodPost:
	port.User.OperationChangePassword(w, r)
case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operations" && p[4] == "change-2fa" && r.Method == http.MethodPost:
	port.User.OperationChangeTwoFactorAuthentication(w, r)
case n == 4 && p[1] == "v1" && p[2] == "users" && p[3] == "select-options" && r.Method == http.MethodGet:
	port.User.ListAsSelectOptions(w, r)

// --- ATTACHMENTS --- //
case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodGet:
	port.Attachment.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodPost:
	port.Attachment.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodGet:
	port.Attachment.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodPut:
	port.Attachment.UpdateByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodDelete:
	port.Attachment.DeleteByID(w, r, p[3])

	// --- PAYMENT PROCESSOR --- //
case n == 5 && p[1] == "v1" && p[2] == "stripe" && p[3] == "create-checkout-session-for-comic-submission" && r.Method == http.MethodPost:
	port.StripePaymentProcessor.CreateStripeCheckoutSessionURLForComicSubmissionID(w, r, p[4])
// case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "complete-checkout-session" && r.Method == http.MethodGet:
// 	port.PaymentProcessor.CompleteStripeCheckoutSession(w, r)
// case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "cancel-subscription" && r.Method == http.MethodPost:
// 	port.PaymentProcessor.CancelStripeSubscription(w, r)
// // case n == 4 && p[1] == "v1" && p[2] == "public" && p[3] == "stripe-webhook":
// // 	port.PaymentProcessor.StripeWebhook(w, r)
// case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "receipts" && r.Method == http.MethodGet:
// 	port.PaymentProcessor.ListLatestStripeReceipts(w, r)
case n == 4 && p[1] == "v1" && p[2] == "public" && p[3] == "stripe-webhook":
	port.StripePaymentProcessor.Webhook(w, r)

// --- OFFERS --- //
case n == 3 && p[1] == "v1" && p[2] == "offers" && r.Method == http.MethodGet:
	port.Offer.List(w, r)
// case n == 3 && p[1] == "v1" && p[2] == "offers" && r.Method == http.MethodPost:
// 	port.Offer.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodGet:
	port.Offer.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodPut:
	port.Offer.UpdateByID(w, r, p[3])
// case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodDelete:
// 	port.Offer.DeleteByID(w, r, p[3])
// case n == 5 && p[1] == "v1" && p[2] == "offer" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
// 	port.Offer.OperationCreateComment(w, r)
case n == 4 && p[1] == "v1" && p[2] == "offers" && p[3] == "select-options" && r.Method == http.MethodGet:
	port.Offer.ListAsSelectOptions(w, r)
case n == 5 && p[1] == "v1" && p[2] == "offer" && p[3] == "service-type" && r.Method == http.MethodGet:
	port.Offer.GetByServiceType(w, r, p[4])

// --- CREDITS --- //
case n == 3 && p[1] == "v1" && p[2] == "credits" && r.Method == http.MethodGet:
	port.Credit.List(w, r)
case n == 3 && p[1] == "v1" && p[2] == "credits" && r.Method == http.MethodPost:
	port.Credit.Create(w, r)
case n == 4 && p[1] == "v1" && p[2] == "credit" && r.Method == http.MethodGet:
	port.Credit.GetByID(w, r, p[3])
case n == 4 && p[1] == "v1" && p[2] == "credit" && r.Method == http.MethodPut:
	port.Credit.UpdateByID(w, r, p[3])

// --- USER PURCHASES --- //
case n == 3 && p[1] == "v1" && p[2] == "user-purchases" && r.Method == http.MethodGet:
	port.UserPurchase.List(w, r)

// --- CATCH ALL: D.N.E. ---
default:
	port.Logger.Debug("404 request",
		slog.Int("n", n),
		slog.String("m", r.Method),
		slog.Any("p", p),
	)
	http.NotFound(w, r)
}
*/
