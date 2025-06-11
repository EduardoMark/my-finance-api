package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/EduardoMark/my-finance-api/pkg/httpResponse"
	"github.com/EduardoMark/my-finance-api/pkg/token"
)

func AuthMiddleware(jwtManager *token.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpResponse.Unauthorized(w)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpResponse.Unauthorized(w)
				return
			}

			token := parts[1]

			claims, err := jwtManager.VerifyToken(token)
			if err != nil {
				httpResponse.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), "name", claims.Name)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "exp", claims.ExpiresAt)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
