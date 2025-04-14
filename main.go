package main

import (
	"avito_test/api/http"
	"avito_test/config"
	"avito_test/pkg"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository/postgreSQL"
	"avito_test/repository/prometheus"
	"avito_test/usecases/service"
	"github.com/go-chi/chi/v5"
	"log"
)

func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	prometheus.InitPrometheus()

	storage, err := postgres_connect.NewPostgresStorage(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed creating Postgres: %s", err.Error())
	}

	UserRepo := postgreSQL.NewUserRepo(storage)
	UserService := service.NewUserService(UserRepo)
	UserHandlers := http.NewUserHandler(UserService)

	PvzRepo := postgreSQL.NewPvzRepo(storage)
	PvzService := service.NewPvzService(PvzRepo)
	PvzHandlers := http.NewPvzHandler(PvzService)

	ReceptionRepo := postgreSQL.NewReceptionRepo(storage)
	ReceptionService := service.NewReceptionService(ReceptionRepo, PvzRepo)
	ReceptionHandlers := http.NewReceptionHandler(ReceptionService)

	ProductRepo := postgreSQL.NewProductRepo(storage)
	ProductService := service.NewProductService(ProductRepo, ReceptionRepo, PvzRepo)
	ProductHandlers := http.NewProductHandler(ProductService)

	r := chi.NewRouter()
	UserHandlers.WithUserHandlers(r)

	r.Use(http.PrometheusMiddleware)

	r.Route("/", func(r chi.Router) {
		r.With(http.AuthMiddleware([]string{"moderator"})).Post("/pvz", PvzHandlers.OpenPvzHandler)
		r.With(http.AuthMiddleware([]string{"employee", "moderator"})).Get("/pvz", PvzHandlers.GetPvzListHandler)
		r.With(http.AuthMiddleware([]string{"employee"})).Group(func(r chi.Router) {
			ReceptionHandlers.WithReceptionHandlers(r)
			ProductHandlers.WithProductHandlers(r)
		})
	})

	log.Printf("Starting server on %s", cfg.Address)
	if err := pkg.CreateAndRunServer(r, cfg.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
