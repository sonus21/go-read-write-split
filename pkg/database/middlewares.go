package database

import (
	"net/http"
)

// Middleware a middleware that is used in chi to set the database context
func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// identify database name from request
			name := dbNameFromRequest(r)
			// set db name in the context
			ctx := SetDb(r.Context(), name)
			// call next handler with database context
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
