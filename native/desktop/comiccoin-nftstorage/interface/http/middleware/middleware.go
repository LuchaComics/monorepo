package middleware

import (
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/blacklist"
)

type Middleware interface {
	Attach(fn http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	Logger    *slog.Logger
	Blacklist blacklist.Provider
}

func NewMiddleware(
	loggerp *slog.Logger,
	blp blacklist.Provider,
) Middleware {
	return &middleware{
		Logger:    loggerp,
		Blacklist: blp,
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) Attach(fn http.HandlerFunc) http.HandlerFunc {
	mid.Logger.Debug("middleware executed")
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `RateLimitMiddleware` will be executed first and
	//     `ProtectedURLsMiddleware` will be executed last.
	fn = mid.EnforceBlacklistMiddleware(fn)
	fn = mid.IPAddressMiddleware(fn)
	fn = mid.URLProcessorMiddleware(fn)
	fn = mid.RateLimitMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}
