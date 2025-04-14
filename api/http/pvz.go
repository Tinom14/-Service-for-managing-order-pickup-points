package http

import (
	"avito_test/api/http/types"
	"avito_test/repository/prometheus"
	"avito_test/usecases"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Pvz struct {
	Service usecases.Pvz
}

func NewPvzHandler(service usecases.Pvz) *Pvz {
	return &Pvz{Service: service}
}

func (p *Pvz) OpenPvzHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateOpenPvzHandlerRequest(r)
	if errors.Is(err, types.ErrInvalidCity) {
		http.Error(w, "Invalid City", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	pvz, err := p.Service.OpenPvz(req.City)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	prometheus.RecordPVZCreated()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pvz); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (p *Pvz) GetPvzListHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateListPvzHandlerRequest(r)
	if err != nil {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	pvzList, err := p.Service.GetPvzListWithFilter(req.StartDate, req.EndDate, req.Page, req.Limit)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pvzList); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}

func (p *Pvz) WithPvzHandlers(r chi.Router) {
	r.Post("/pvz", p.OpenPvzHandler)
	r.Get("/pvz", p.OpenPvzHandler)
}
