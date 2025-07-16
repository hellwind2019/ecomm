package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hellwind2019/ecomm/token"
)

type authKey struct{}

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
func GetAdminMiddlewareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
				return
			}
			if !claims.IsAdmin {
				http.Error(w, "user is not an admin", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}
}
func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization token is missing")
	}
	fiels := strings.Fields(authHeader)
	if len(fiels) != 2 || fiels[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}
	token := fiels[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return claims, nil
}
