package middleware

import (
	"net/http"
)

func (mid *middleware) PostJWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		//TODO: Impl.
		fn(w, r.WithContext(ctx)) // Flow to the next middleware.
	}
}
