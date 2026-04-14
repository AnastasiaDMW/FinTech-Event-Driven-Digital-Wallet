package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIdKey contextKey = "userId"

type AuthMiddleware struct {
	PublicKey *rsa.PublicKey
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.PublicKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		rawUserID, ok := claims["userId"]
		if !ok {
			http.Error(w, "userId missing", http.StatusUnauthorized)
			return
		}

		var userID int64

		switch v := rawUserID.(type) {
		case float64:
			userID = int64(v)
		case string:
			parsed, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				http.Error(w, "invalid userId format", http.StatusUnauthorized)
				return
			}
			userID = parsed
		default:
			http.Error(w, "invalid userId type", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}