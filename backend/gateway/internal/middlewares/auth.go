package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"gateway/internal/helper"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

type JWTClaims struct {
	ID string `json:"id"`
	jwtlib.RegisteredClaims
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		httpInfo := map[string]string{
			"method": r.Method,
			"url":    r.URL.String(),
		}

		authorization := r.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			durationMs := time.Since(start).Milliseconds()
			helper.Response(helper.Response_struct{Error: "missing bearer token"}, w, http.StatusUnauthorized, "token required", httpInfo, durationMs)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		if tokenString == "" {
			durationMs := time.Since(start).Milliseconds()
			helper.Response(helper.Response_struct{Error: "empty bearer token"}, w, http.StatusUnauthorized, "token required", httpInfo, durationMs)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			durationMs := time.Since(start).Milliseconds()
			helper.Response(helper.Response_struct{Error: "jwt secret not configured"}, w, http.StatusInternalServerError, "internal server error", httpInfo, durationMs)
			return
		}

		claims := &JWTClaims{}
		token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			durationMs := time.Since(start).Milliseconds()
			helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusUnauthorized, "invalid token", httpInfo, durationMs)
			return
		}

		if claims.ID == "" {
			durationMs := time.Since(start).Milliseconds()
			helper.Response(helper.Response_struct{Error: "missing user id in token"}, w, http.StatusUnauthorized, "invalid token", httpInfo, durationMs)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDContextKey).(string)
	return userID
}
