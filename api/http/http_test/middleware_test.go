package http_test

import (
	http2 "avito_test/api/http"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	validToken, err := generateTestToken("1", "employee")
	if err != nil {
		t.Fatal(err)
	}

	r := chi.NewRouter()
	r.With(http2.AuthMiddleware([]string{"employee"})).
		Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func generateTestToken(userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte("secret-key"))
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := chi.NewRouter()
	r.With(http2.AuthMiddleware([]string{"employee"})).
		Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
