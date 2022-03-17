package api

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

// MiddlewareAuth gets token from request, checks it and sets user_id to context
func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if token == nil {
			http.Error(w, "token is nil", http.StatusUnauthorized)
			return
		}

		if err := jwt.Validate(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Token is authenticated, parse claims
		userID, ok := claims["user_id"]
		if !ok {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Set userID to context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
