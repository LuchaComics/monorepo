package middleware

import (
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/blacklist"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type Middleware interface {
	Attach(fn http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	logger                    *slog.Logger
	blacklist                 blacklist.Provider
	jwt                       jwt.Provider
	userGetBySessionIDUseCase *usecase.UserGetBySessionIDUseCase
}

func NewMiddleware(
	loggerp *slog.Logger,
	blp blacklist.Provider,
	jwtp jwt.Provider,
	uc1 *usecase.UserGetBySessionIDUseCase,
) Middleware {
	return &middleware{
		logger:                    loggerp,
		blacklist:                 blp,
		jwt:                       jwtp,
		userGetBySessionIDUseCase: uc1,
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) Attach(fn http.HandlerFunc) http.HandlerFunc {
	mid.logger.Debug("middleware executed")
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `RateLimitMiddleware` will be executed first and
	//     `ProtectedURLsMiddleware` will be executed last.
	// fn = mid.ProtectedURLsMiddleware(fn)
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
