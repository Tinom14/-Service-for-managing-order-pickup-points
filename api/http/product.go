package http

import (
	"avito_test/api/http/types"
	"avito_test/repository"
	"avito_test/repository/prometheus"
	"avito_test/usecases"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Product struct {
	Service usecases.Product
}

func NewProductHandler(service usecases.Product) *Product {
	return &Product{Service: service}
}

func (p *Product) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateAddProductHandlerRequest(r)
	if errors.Is(err, types.ErrTypePvzIdRequired) {
		http.Error(w, "Type and pvzId are required", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	pvzId, err := strconv.Atoi(req.PvzId)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	product, err := p.Service.AddProduct(req.Type, pvzId)
	if err != nil {
		switch {
		case errors.Is(err, repository.NotFound):
			http.Error(w, "Pvz not found", http.StatusBadRequest)
		case errors.Is(err, usecases.ErrAlreadyClosed):
			http.Error(w, "Reception already closed", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}
	prometheus.RecordProductAdded()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (p *Product) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzId := chi.URLParam(r, "pvzId")

	pvzIdInt, err := strconv.Atoi(pvzId)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = p.Service.DeleteProduct(pvzIdInt)
	if err != nil {
		switch {
		case errors.Is(err, repository.NotFound):
			http.Error(w, "Pvz not found", http.StatusBadRequest)
		case errors.Is(err, usecases.ErrAlreadyClosed):
			http.Error(w, "Reception already closed", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Product) WithProductHandlers(r chi.Router) {
	r.Post("/products", p.AddProductHandler)
	r.Post("/pvz/{pvzId}/delete_last_product", p.DeleteProductHandler)
}
