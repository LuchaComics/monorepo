package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/ratelimit"

	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/gateway/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/blacklist"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/time"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/uuid"
)

type Middleware interface {
	Attach(fn http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	Config            *config.Conf
	Logger            *slog.Logger
	Time              time.Provider
	JWT               jwt.Provider
	UUID              uuid.Provider
	Blacklist         blacklist.Provider
	GatewayController gateway_c.GatewayController
}

func NewMiddleware(
	configp *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	timep time.Provider,
	jwtp jwt.Provider,
	blp blacklist.Provider,
	gatewayController gateway_c.GatewayController,
) Middleware {
	return &middleware{
		Logger:            loggerp,
		UUID:              uuidp,
		Time:              timep,
		JWT:               jwtp,
		Blacklist:         blp,
		GatewayController: gatewayController,
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) Attach(fn http.HandlerFunc) http.HandlerFunc {
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `RateLimitMiddleware` will be executed first and
	//     `ProtectedURLsMiddleware` will be executed last.
	fn = mid.ProtectedURLsMiddleware(fn)
	fn = mid.PostJWTProcessorMiddleware(fn) // Note: Must be above `JWTProcessorMiddleware`.
	fn = mid.JWTProcessorMiddleware(fn)     // Note: Must be above `PreJWTProcessorMiddleware`.
	fn = mid.PreJWTProcessorMiddleware(fn)  // Note: Must be above `URLProcessorMiddleware` and `IPAddressMiddleware`.
	fn = mid.EnforceBlacklistMiddleware(fn)
	fn = mid.IPAddressMiddleware(fn)
	fn = mid.URLProcessorMiddleware(fn)
	fn = mid.RateLimitMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}

func (mid *middleware) RateLimitMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Special thanks: https://ubogdan.com/2021/09/ip-based-rate-limit-middleware-using-go.uber.org/ratelimit/
	var lmap sync.Map

	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			mid.Logger.Error("invalid RemoteAddr", slog.Any("err", err), slog.Any("middleware", "RateLimitMiddleware"))
			http.Error(w, fmt.Sprintf("invalid RemoteAddr: %s", err), http.StatusInternalServerError)
			return
		}

		lif, ok := lmap.Load(host)
		if !ok {
			lif = ratelimit.New(50) // per second.
		}

		lm, ok := lif.(ratelimit.Limiter)
		if !ok {
			mid.Logger.Error("internal middleware error: typecast failed", slog.Any("middleware", "RateLimitMiddleware"))
			http.Error(w, "internal middleware error: typecast failed", http.StatusInternalServerError)
			return
		}

		lm.Take()
		lmap.Store(host, lm)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

// Note: This middleware must have `IPAddressMiddleware` executed first before running.
func (mid *middleware) EnforceBlacklistMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)
		proxies, _ := ctx.Value(constants.SessionProxies).(string)

		// Case 1 of 2: Check banned IP addresses.
		if mid.Blacklist.IsBannedIPAddress(ipAddress) {

			// If the client IP address is banned, check to see if the client
			// is making a call to a URL which is not banned, and if the URL
			// is not banned (has not been banned before) then print it to
			// the console logs for future analysis. Else if the URL is banned
			// then don't bother printing to console. The purpose of this code
			// is to not clog the console log with warnings.
			if !mid.Blacklist.IsBannedURL(r.URL.Path) {
				mid.Logger.Warn("rejected request by ip",
					slog.Any("url", r.URL.Path),
					slog.String("ip_address", ipAddress),
					slog.String("proxies", proxies),
					slog.Any("middleware", "EnforceBlacklistMiddleware"))
			}
			http.Error(w, "forbidden at this time", http.StatusForbidden)
			return
		}

		// Case 2 of 2: Check banned URL.
		if mid.Blacklist.IsBannedURL(r.URL.Path) {

			// If the URL is banned, check to see if the client IP address is
			// banned, and if the client has not been banned before then print
			// to console the new offending client ip. Else do not print if
			// the offending client IP address has been banned before. The
			// purpose of this code is to not clog the console log with warnings.
			if !mid.Blacklist.IsBannedIPAddress(ipAddress) {
				mid.Logger.Warn("rejected request by url",
					slog.Any("url", r.URL.Path),
					slog.String("ip_address", ipAddress),
					slog.String("proxies", proxies),
					slog.Any("middleware", "EnforceBlacklistMiddleware"))
			}

			// DEVELOPERS NOTE:
			// Simply return a 404, but in our console log we can see the IP
			// address whom made this call.
			http.Error(w, "does not exist at this time", http.StatusNotFound)
			return
		}

		next(w, r.WithContext(ctx))
	}
}

// URLProcessorMiddleware Middleware will split the full URL path into slash-sperated parts and save to
// the context to flow downstream in the app for this particular request.
func (mid *middleware) URLProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Split path into slash-separated parts, for example, path "/foo/bar"
		// gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
		// "/api", as a result we will start the array slice at "1".
		p := strings.Split(r.URL.Path, "/")[1:]

		// log.Println(p) // For debugging purposes only.

		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "url_split", p)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

// PreJWTProcessorMiddleware checks to see if we are visiting an unprotected URL and if so then
// let the system know we need to skip authorization handling.
func (mid *middleware) PreJWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"health-check":    true,
			"version":         true,
			"greeting":        true,
			"login":           true,
			"refresh-token":   true,
			"register":        true,
			"verify":          true,
			"forgot-password": true,
			"password-reset":  true,
			"cpsrn":           true,
			"select-options":  true,
			"public":          true,
			"ipfs":            true,
		}

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size then do not skip authorization.
		if len(urlSplit) < 3 {
			// mid.Logger.Warn("Skipcps-ipfsstoreg authorization | len less then 3")
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, false)
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// mid.Logger.Debug("api call",
		// 	slog.Any("url", urlSplit),
		// 	slog.Any("urlSplit[2]", urlSplit[2]),
		// 	slog.Any("middleware", "PreJWTProcessorMiddleware"))

		// Skip authorization if the URL matches the whitelist else we need to
		// run authorization check.
		if skipPath[urlSplit[2]] {
			// mid.Logger.Warn("Skipcps-ipfsstoreg authorization | skipPath found")
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, true)
		} else {
			// For debugging purposes only.
			// log.Println("PreJWTProcessorMiddleware | Protected URL detected")
			// log.Println("PreJWTProcessorMiddleware | urlSplit:", urlSplit)
			// log.Println("PreJWTProcessorMiddleware | urlSplit[2]:", urlSplit[2])
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, false)
		}

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) JWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipcps-ipfsstoreg authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// Extract our auth header array.
		reqToken := r.Header.Get("Authorization")

		// For debugging purposes.
		// log.Println("JWTProcessorMiddleware | reqToken:", reqToken)

		// Before running our JWT middleware we need to confirm there is an
		// an `Authorization` header to run our middleware. This is an important
		// step!
		if reqToken != "" && strings.Contains(reqToken, "undefined") == false {

			// Special thanks to "poise" via https://stackoverflow.com/a/44700761
			splitToken := strings.Split(reqToken, "JWT ")
			if len(splitToken) < 2 {
				mid.Logger.Warn("not properly formatted authorization header", slog.Any("middleware", "JWTProcessorMiddleware"))
				http.Error(w, "not properly formatted authorization header", http.StatusBadRequest)
				return
			}

			reqToken = splitToken[1]
			// log.Println("JWTProcessorMiddleware | reqToken:", reqToken) // For debugging purposes only.

			sessionID, err := mid.JWT.ProcessJWTToken(reqToken)
			// log.Println("JWTProcessorMiddleware | sessionUUID:", sessionUUID) // For debugging purposes only.

			if err == nil {
				// Update our context to save our JWT token content information.
				ctx = context.WithValue(ctx, constants.SessionIsAuthorized, true)
				ctx = context.WithValue(ctx, constants.SessionID, sessionID)

				// Flow to the next middleware with our JWT token saved.
				fn(w, r.WithContext(ctx))
				return
			}

			// The following code will lookup the URL path in a whitelist and
			// if the visited path matches then we will skip any token errors.
			// We do this because a majority of API endpoints are protected
			// by authorization.

			urlSplit := ctx.Value("url_split").([]string)
			skipPath := map[string]bool{
				"health-check":    true,
				"version":         true,
				"greeting":        true,
				"login":           true,
				"refresh-token":   true,
				"register":        true,
				"verify":          true,
				"forgot-password": true,
				"password-reset":  true,
				"cpsrn":           true,
				"select-options":  true,
				"public":          true,
				"ipfs":            true,
			}

			// DEVELOPERS NOTE:
			// If the URL cannot be split into the size we want then skip running
			// this middleware.
			if len(urlSplit) >= 3 {
				if skipPath[urlSplit[2]] {
					mid.Logger.Warn("Skipcps-ipfsstoreg expired or error token", slog.Any("middleware", "JWTProcessorMiddleware"))
				} else {
					// For debugging purposes only.
					// log.Println("JWTProcessorMiddleware | ProcessJWT | err", err, "for reqToken:", reqToken)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit:", urlSplit)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit[2]:", urlSplit[2])
					mid.Logger.Warn("unauthorized api call", slog.Any("url", urlSplit), slog.Any("middleware", "JWTProcessorMiddleware"))
					http.Error(w, "attempting to access a protected endpoint", http.StatusUnauthorized)
					return
				}
			}
		}

		// Flow to the next middleware without anything done.
		ctx = context.WithValue(ctx, constants.SessionIsAuthorized, false)
		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) PostJWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip this middleware if user is on a whitelisted URL path.
		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipcps-ipfsstoreg authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// urlSplit := ctx.Value("url_split").([]string)
		// mid.Logger.Debug("api call",
		// 	slog.Any("url", urlSplit),
		// 	slog.Any("skipAuthorization", skipAuthorization),
		// 	slog.Any("middleware", "PostJWTProcessorMiddleware"))

		// Get our authorization information.
		isAuthorized, ok := ctx.Value(constants.SessionIsAuthorized).(bool)
		if ok && isAuthorized {
			sessionID := ctx.Value(constants.SessionID).(string)

			// Lookup our user profile in the session or return 500 error.
			user, err := mid.GatewayController.GetUserBySessionID(ctx, sessionID) //TODO: IMPLEMENT.
			if err != nil {
				mid.Logger.Warn("GetUserBySessionID error", slog.Any("err", err), slog.Any("middleware", "PostJWTProcessorMiddleware"))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If no user was found then that means our session expired and the
			// user needs to login or use the refresh token.
			if user == nil {
				mid.Logger.Warn("Session expired - please log in again", slog.Any("middleware", "PostJWTProcessorMiddleware"))
				http.Error(w, "attempting to access a protected endpoint", http.StatusUnauthorized)
				return
			}

			// // If system administrator disabled the user account then we need
			// // to generate a 403 error letting the user know their account has
			// // been disabled and you cannot access the protected API endpoint.
			// if user.State == 0 {
			// 	http.Error(w, "Account disabled - please contact admin", http.StatusForbidden)
			// 	return
			// }

			// Save our user information to the context.
			// Save our user.
			ctx = context.WithValue(ctx, constants.SessionUser, user)

			// // For debugging purposes only.
			// mid.Logger.Debug("Fetched session record",
			// 	slog.Any("ID", user.ID),
			// 	slog.String("SessionID", sessionID),
			// 	slog.String("Name", user.Name),
			// 	slog.String("FirstName", user.FirstName),
			// 	slog.String("Email", user.Email))

			// Save individual pieces of the user profile.
			ctx = context.WithValue(ctx, constants.SessionID, sessionID)
			ctx = context.WithValue(ctx, constants.SessionUserID, user.ID)
			ctx = context.WithValue(ctx, constants.SessionUserRole, user.Role)
			ctx = context.WithValue(ctx, constants.SessionUserName, user.Name)
			ctx = context.WithValue(ctx, constants.SessionUserFirstName, user.FirstName)
			ctx = context.WithValue(ctx, constants.SessionUserLastName, user.LastName)
			ctx = context.WithValue(ctx, constants.SessionUserTenantID, user.TenantID)
			ctx = context.WithValue(ctx, constants.SessionUserTenantName, user.TenantName)
			ctx = context.WithValue(ctx, constants.SessionUserTenantLevel, user.TenantLevel)
			ctx = context.WithValue(ctx, constants.SessionUserTenantTimezone, user.TenantTimezone)
		}

		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) IPAddressMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the IPAddress. Code taken from: https://stackoverflow.com/a/55738279
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		// Save our IP address to the context.
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.SessionIPAddress, IPAddress)
		fn(w, r.WithContext(ctx)) // Flow to the next middleware.
	}
}

// ProtectedURLsMiddleware The purpose of this middleware is to return a `401 unauthorized` error if
// the user is not authorized when visiting a protected URL.
func (mid *middleware) ProtectedURLsMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip this middleware if user is on a whitelisted URL path.
		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipcps-ipfsstoreg authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"health-check":    true,
			"version":         true,
			"greeting":        true,
			"login":           true,
			"refresh-token":   true,
			"verify":          true,
			"forgot-password": true,
			"password-reset":  true,
			"cpsrn":           true,
			"select-options":  true,
			"public":          true,
			"ipfs":            true,
		}

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size we want then skip running
		// this middleware.
		if len(urlSplit) < 3 {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// mid.Logger.Debug("api call",
		// 	slog.Any("url", urlSplit),
		// 	slog.Any("urlSplit[2]", urlSplit[2]),
		// 	slog.Any("middleware", "ProtectedURLsMiddleware"))

		if skipPath[urlSplit[2]] {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
		} else {
			// Get our authorization information.
			isAuthorized, ok := ctx.Value(constants.SessionIsAuthorized).(bool)

			// Either accept continuing execution or return 401 error.
			if ok && isAuthorized {
				fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			} else {
				mid.Logger.Warn("unauthorized api call", slog.Any("url", urlSplit), slog.Any("middleware", "ProtectedURLsMiddleware"))
				http.Error(w, "attempting to access a protected endpoint", http.StatusUnauthorized)
				return
			}
		}
	}
}
