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

type Reception struct {
	Service usecases.Reception
}

func NewReceptionHandler(service usecases.Reception) *Reception {
	return &Reception{Service: service}
}

func (rec *Reception) StartReceptionHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateStartReceptionHandlerRequest(r)
	if err != nil {
		http.Error(w, "pvzId is required", http.StatusBadRequest)
		return
	}

	pvzId, err := strconv.Atoi(req.PvzId)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	reception, err := rec.Service.StartReception(pvzId)
	if err != nil {
		switch {
		case errors.Is(err, repository.NotFound):
			http.Error(w, "Pvz not found", http.StatusBadRequest)
		case errors.Is(err, types.ErrPvzIdRequired):
			http.Error(w, "PvzId is required", http.StatusBadRequest)
		case errors.Is(err, usecases.ErrUnclosedReception):
			http.Error(w, "Unclosed reception", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reception); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (rec *Reception) CloseReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzId := chi.URLParam(r, "pvzId")

	pvzIdInt, err := strconv.Atoi(pvzId)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	reception, err := rec.Service.CloseReception(pvzIdInt)
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
	prometheus.RecordReceptionCreated()

	if err := json.NewEncoder(w).Encode(reception); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (rec *Reception) WithReceptionHandlers(r chi.Router) {
	r.Post("/receptions", rec.StartReceptionHandler)
	r.Post("/pvz/{pvzId}/close_last_reception", rec.CloseReceptionHandler)
}
