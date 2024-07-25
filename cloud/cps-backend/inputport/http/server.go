package http

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/rs/cors"

	attachment "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/httptransport"
	comicsub "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/httptransport"
	credit "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/httptransport"
	customer "github.com/LuchaComics/monorepo/cloud/cps-backend/app/customer/httptransport"
	gateway "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/httptransport"
	offer "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/httptransport"
	strpp "github.com/LuchaComics/monorepo/cloud/cps-backend/app/paymentprocessor/httptransport/stripe"
	receipt "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/httptransport"
	store "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/httptransport"
	user "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/httptransport"
	userpurchase "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/inputport/http/middleware"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type httpInputPort struct {
	Config                 *config.Conf
	Logger                 *slog.Logger
	Server                 *http.Server
	Middleware             middleware.Middleware
	Gateway                *gateway.Handler
	User                   *user.Handler
	Store                  *store.Handler
	ComicSubmission        *comicsub.Handler
	Customer               *customer.Handler
	Attachment             *attachment.Handler
	Offer                  *offer.Handler
	Receipt                *receipt.Handler
	UserPurchase           *userpurchase.Handler
	StripePaymentProcessor *strpp.Handler
	Credit                 *credit.Handler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	mid middleware.Middleware,
	gh *gateway.Handler,
	cu *user.Handler,
	org *store.Handler,
	t *comicsub.Handler,
	cust *customer.Handler,
	att *attachment.Handler,
	off *offer.Handler,
	inv *receipt.Handler,
	usrp *userpurchase.Handler,
	strpp *strpp.Handler,
	cr *credit.Handler,
) InputPortServer {
	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	addr := fmt.Sprintf("%s:%s", configp.AppServer.IP, configp.AppServer.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Create our HTTP server controller.
	p := &httpInputPort{
		Config:                 configp,
		Logger:                 loggerp,
		Middleware:             mid,
		Gateway:                gh,
		User:                   cu,
		Store:                  org,
		Offer:                  off,
		ComicSubmission:        t,
		Customer:               cust,
		Attachment:             att,
		Receipt:                inv,
		UserPurchase:           usrp,
		StripePaymentProcessor: strpp,
		Credit:                 cr,
		Server:                 srv,
	}

	// Attach the HTTP server controller to the ServerMux.
	mux.HandleFunc("/", mid.Attach(p.HandleRequests))

	return p
}

func (port *httpInputPort) Run() {
	port.Logger.Info("HTTP server running")
	if err := port.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		port.Logger.Error("listen failed", slog.Any("error", err))

		// DEVELOPERS NOTE: We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of app.
		panic("failed running")
	}
}

func (port *httpInputPort) Shutdown() {
	port.Logger.Info("HTTP server shutdown")
}

func (port *httpInputPort) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)

	switch {
	// --- GATEWAY & PROFILE & DASHBOARD --- //
	case n == 3 && p[1] == "v1" && p[2] == "health-check" && r.Method == http.MethodGet:
		port.Gateway.HealthCheck(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "version" && r.Method == http.MethodGet:
		port.Gateway.Version(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "greeting" && r.Method == http.MethodPost:
		port.Gateway.Greet(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "login" && r.Method == http.MethodPost:
		port.Gateway.Login(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "register" && p[3] == "business" && r.Method == http.MethodPost:
		port.Gateway.RegisterBusiness(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "register" && p[3] == "customer" && r.Method == http.MethodPost:
		port.Gateway.RegisterCustomer(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "refresh-token" && r.Method == http.MethodPost:
		port.Gateway.RefreshToken(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "verify" && r.Method == http.MethodPost:
		port.Gateway.Verify(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "logout" && r.Method == http.MethodPost:
		port.Gateway.Logout(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodGet:
		port.Gateway.Profile(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodPut:
		port.Gateway.ProfileUpdate(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "profile" && p[3] == "change-password" && r.Method == http.MethodPut:
		port.Gateway.ProfileChangePassword(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "forgot-password" && r.Method == http.MethodPost:
		port.Gateway.ForgotPassword(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "password-reset" && r.Method == http.MethodPost:
		port.Gateway.PasswordReset(w, r)
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
}
