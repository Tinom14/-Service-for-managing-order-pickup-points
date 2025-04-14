//go:build integration

package integration_test

import (
	http2 "avito_test/api/http"
	"avito_test/config"
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository/postgreSQL"
	"avito_test/usecases"
	"avito_test/usecases/service"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *chi.Mux
	token  string
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	cfg := config.AppConfig{
		HTTPConfig: config.HTTPConfig{
			Address: ":8080",
		},
		Postgres: config.Postgres{
			Host:     os.Getenv("TEST_DB_HOST"),
			Port:     5432,
			User:     os.Getenv("TEST_DB_USER"),
			Password: os.Getenv("TEST_DB_PASSWORD"),
			DBName:   os.Getenv("TEST_DB_NAME"),
			SSLMode:  "disable",
		},
	}

	storage, err := postgres_connect.NewPostgresStorage(cfg.Postgres)
	if err != nil {
		s.T().Fatalf("failed to connect to test database: %s", err)
	}

	s.cleanDatabase(storage)

	userRepo := postgreSQL.NewUserRepo(storage)
	pvzRepo := postgreSQL.NewPvzRepo(storage)
	receptionRepo := postgreSQL.NewReceptionRepo(storage)
	productRepo := postgreSQL.NewProductRepo(storage)

	userService := service.NewUserService(userRepo)
	pvzService := service.NewPvzService(pvzRepo)
	receptionService := service.NewReceptionService(receptionRepo, pvzRepo)
	productService := service.NewProductService(productRepo, receptionRepo, pvzRepo)

	userHandler := http2.NewUserHandler(userService)
	pvzHandler := http2.NewPvzHandler(pvzService)
	receptionHandler := http2.NewReceptionHandler(receptionService)
	productHandler := http2.NewProductHandler(productService)

	s.token = s.createTestUserAndGetToken(userService)

	s.router = chi.NewRouter()
	userHandler.WithUserHandlers(s.router)

	s.router.Route("/", func(r chi.Router) {
		r.With(http2.AuthMiddleware([]string{"moderator"})).Post("/pvz", pvzHandler.OpenPvzHandler)
		r.With(http2.AuthMiddleware([]string{"employee", "moderator"})).Get("/pvz", pvzHandler.GetPvzListHandler)
		r.With(http2.AuthMiddleware([]string{"employee"})).Group(func(r chi.Router) {
			receptionHandler.WithReceptionHandlers(r)
			productHandler.WithProductHandlers(r)
		})
	})
}

func (s *IntegrationTestSuite) cleanDatabase(storage *postgres_connect.PostgresStorage) {
	_, err := storage.Db.Exec(`
		TRUNCATE TABLE users, pvz, receptions, products, reception_products RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		s.T().Fatalf("failed to clean test database: %s", err)
	}
}

func (s *IntegrationTestSuite) createTestUserAndGetToken(service usecases.User) string {
	_, err := service.Register("moderator@test.com", "password123", "moderator")
	if err != nil {
		s.T().Fatalf("failed to create test user: %s", err)
	}

	token, err := service.Login("moderator@test.com", "password123")
	if err != nil {
		s.T().Fatalf("failed to get test token: %s", err)
	}

	return token
}

func (s *IntegrationTestSuite) TestFullPvzWorkflow() {
	pvz := s.createPvz("Москва")
	assert.NotEmpty(s.T(), pvz.Id)

	reception := s.startReception(pvz.Id)
	assert.NotEmpty(s.T(), reception.Id)
	assert.Equal(s.T(), "in_progress", reception.Status)

	for i := 0; i < 5; i++ {
		product := s.addProduct("product_"+strconv.Itoa(i+1), pvz.Id)
		assert.NotEmpty(s.T(), product.Id)
	}

	closedReception := s.closeReception(pvz.Id)
	assert.Equal(s.T(), "closed", closedReception.Status)
}

func (s *IntegrationTestSuite) createPvz(city string) domain.Pvz {
	reqBody := struct {
		City string `json:"city"`
	}{City: city}

	var resp domain.Pvz
	s.sendAuthenticatedRequest("POST", "/pvz", reqBody, &resp, "moderator")
	return resp
}

func (s *IntegrationTestSuite) startReception(pvzId int) domain.Reception {
	reqBody := struct {
		PvzId string `json:"pvzId"`
	}{PvzId: strconv.Itoa(pvzId)}

	var resp domain.Reception
	s.sendAuthenticatedRequest("POST", "/receptions", reqBody, &resp, "employee")
	return resp
}

func (s *IntegrationTestSuite) addProduct(productType string, pvzId int) domain.Product {
	reqBody := struct {
		Type  string `json:"type"`
		PvzId string `json:"pvzId"`
	}{
		Type:  productType,
		PvzId: strconv.Itoa(pvzId),
	}

	var resp domain.Product
	s.sendAuthenticatedRequest("POST", "/products", reqBody, &resp, "employee")
	return resp
}

func (s *IntegrationTestSuite) closeReception(pvzId int) domain.Reception {
	var resp domain.Reception
	s.sendAuthenticatedRequest("POST", "/pvz/"+strconv.Itoa(pvzId)+"/close_last_reception", nil, &resp, "employee")
	return resp
}

func (s *IntegrationTestSuite) sendAuthenticatedRequest(method, path string, requestBody interface{}, response interface{}, requiredRole string) {
	body, _ := json.Marshal(requestBody)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+s.token)

	s.router.ServeHTTP(w, r)

	assert.Equal(s.T(), http.StatusOK, w.Code, "unexpected status code for %s %s", method, path)

	if w.Body.Len() > 0 {
		err := json.Unmarshal(w.Body.Bytes(), response)
		if err != nil {
			s.T().Fatalf("failed to unmarshal response: %s", err)
		}
	}
}
