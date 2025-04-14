package http

import (
	"avito_test/config"
	"avito_test/repository/prometheus"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(requiredRoles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			id, idOk := claims["id"].(string)
			role, roleOk := claims["role"].(string)
			if !idOk || !roleOk {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if len(requiredRoles) > 0 && !roleAllowed(role, requiredRoles) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "user_info", map[string]interface{}{
				"id":   id,
				"role": role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func roleAllowed(role string, allowedRoles []string) bool {
	for _, r := range allowedRoles {
		if r == role {
			return true
		}
	}
	return false
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseRecorder{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		prometheus.RecordHTTPRequest(r.Method, r.URL.Path, rw.status, duration)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
