package http

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/rs/cors"

	gateway "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/httptransport"
	ipfsgateway "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/httptransport"
	nft "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/httptransport"
	nftasset "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/httptransport"
	collection "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/httptransport"
	tenant "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/httptransport"
	user "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/inputport/http/middleware"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type httpInputPort struct {
	Config      *config.Conf
	Logger      *slog.Logger
	Server      *http.Server
	Middleware  middleware.Middleware
	Gateway     *gateway.Handler
	User        *user.Handler
	Tenant      *tenant.Handler
	Collection  *collection.Handler
	NFT         *nft.Handler
	NFTAsset    *nftasset.Handler
	IpfsGateway *ipfsgateway.Handler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	mid middleware.Middleware,
	gh *gateway.Handler,
	cu *user.Handler,
	org *tenant.Handler,
	co *collection.Handler,
	nft *nft.Handler,
	nftasset *nftasset.Handler,
	ipfsgate *ipfsgateway.Handler,
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
		Config:      configp,
		Logger:      loggerp,
		Middleware:  mid,
		Gateway:     gh,
		User:        cu,
		Tenant:      org,
		Collection:  co,
		NFT:         nft,
		NFTAsset:    nftasset,
		IpfsGateway: ipfsgate,
		Server:      srv,
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
	port.Logger.Info("HTTP server shutting down now...")
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

	// --- ORGANIZATION --- //
	case n == 3 && p[1] == "v1" && p[2] == "tenants" && r.Method == http.MethodGet:
		port.Tenant.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "tenants" && r.Method == http.MethodPost:
		port.Tenant.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "tenant" && r.Method == http.MethodGet:
		port.Tenant.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "tenant" && r.Method == http.MethodPut:
		port.Tenant.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "tenant" && r.Method == http.MethodDelete:
		port.Tenant.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "tenants" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.Tenant.OperationCreateComment(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "tenants" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.Tenant.ListAsSelectOptionByFilter(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "public" && p[3] == "tenants-select-options" && r.Method == http.MethodGet:
		port.Tenant.PublicListAsSelectOptionByFilter(w, r)

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

	// --- COLLECTIONS --- //
	case n == 3 && p[1] == "v1" && p[2] == "nft-collections" && r.Method == http.MethodGet:
		port.Collection.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "nft-collections" && r.Method == http.MethodPost:
		port.Collection.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "nft-collection" && r.Method == http.MethodGet:
		port.Collection.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nft-collection" && r.Method == http.MethodPut:
		port.Collection.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nft-collection" && r.Method == http.MethodDelete:
		port.Collection.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "wallet-balance" && r.Method == http.MethodGet:
		port.Collection.OperationGetWalletBalanceByID(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "deploy-smart-contract" && r.Method == http.MethodPost:
		port.Collection.OperationDeploySmartContract(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "get-token-uri" && r.Method == http.MethodGet:
		port.Collection.OperationGetTokenURI(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "mint" && r.Method == http.MethodPost:
		port.Collection.OperationMint(w, r)
	case n == 6 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "backup" && p[5] == "json" && r.Method == http.MethodPost:
		port.Collection.OperationBackupInJSON(w, r)
	case n == 6 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "backup" && p[5] == "xml" && r.Method == http.MethodPost:
		port.Collection.OperationBackupInXML(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "nft-collections" && p[3] == "operations" && p[4] == "restore" && r.Method == http.MethodPost:
		port.Collection.OperationRestore(w, r)

	// --- NFT METADATA --- //
	case n == 3 && p[1] == "v1" && p[2] == "nfts" && r.Method == http.MethodGet:
		port.NFT.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "nfts" && r.Method == http.MethodPost:
		port.NFT.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "nft" && r.Method == http.MethodGet:
		port.NFT.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nft" && r.Method == http.MethodPut:
		port.NFT.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nft" && r.Method == http.MethodDelete:
		port.NFT.DeleteByID(w, r, p[3])

	// 	// --- NFT ASSETS --- //
	// case n == 3 && p[1] == "v1" && p[2] == "nft-assets" && r.Method == http.MethodGet:
	// 	port.NFT.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "nft-assets" && r.Method == http.MethodPost:
		port.NFTAsset.Create(w, r)
	// case n == 4 && p[1] == "v1" && p[2] == "nft-assets" && r.Method == http.MethodGet:
	// 	port.NFT.GetByID(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "nft-asset" && r.Method == http.MethodPut:
	// 	port.NFT.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nft-asset" && r.Method == http.MethodDelete:
		port.NFTAsset.DeleteByID(w, r, p[3])

	// --- IPFS GATEWAY --- //
	case n == 2 && p[0] == "ipfs" && r.Method == http.MethodGet:
		port.IpfsGateway.GetByCID(w, r, p[1])
	case n == 3 && p[0] == "ipns" && r.Method == http.MethodGet:
		port.IpfsGateway.GetByIPNSPath(w, r, p[1], p[2])

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