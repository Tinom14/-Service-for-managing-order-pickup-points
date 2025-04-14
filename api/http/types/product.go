package types

import (
	"encoding/json"
	"net/http"
)

type AddProductHandlerRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

func CreateAddProductHandlerRequest(r *http.Request) (*AddProductHandlerRequest, error) {
	var req AddProductHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.Type == "" || req.PvzId == "" {
		return nil, ErrTypePvzIdRequired
	}
	return &req, nil
}
