package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
)

func (mid *middleware) JWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.logger.Warn("Skipping authorization")
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
				mid.logger.Warn("not properly formatted authorization header", slog.Any("middleware", "JWTProcessorMiddleware"))
				http.Error(w, "not properly formatted authorization header", http.StatusBadRequest)
				return
			}

			reqToken = splitToken[1]
			// log.Println("JWTProcessorMiddleware | reqToken:", reqToken) // For debugging purposes only.

			sessionID, err := mid.jwt.ProcessJWTToken(reqToken)
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
			}

			// DEVELOPERS NOTE:
			// If the URL cannot be split into the size we want then skip running
			// this middleware.
			if len(urlSplit) >= 3 {
				if skipPath[urlSplit[2]] {
					mid.logger.Warn("Skipping expired or error token", slog.Any("middleware", "JWTProcessorMiddleware"))
				} else {
					// For debugging purposes only.
					// log.Println("JWTProcessorMiddleware | ProcessJWT | err", err, "for reqToken:", reqToken)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit:", urlSplit)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit[2]:", urlSplit[2])
					mid.logger.Warn("unauthorized api call", slog.Any("url", urlSplit), slog.Any("middleware", "JWTProcessorMiddleware"))
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
