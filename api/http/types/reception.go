package types

import (
	"encoding/json"
	"net/http"
)

type StartReceptionHandlerRequest struct {
	PvzId string `json:"pvzId"`
}

func CreateStartReceptionHandlerRequest(r *http.Request) (*StartReceptionHandlerRequest, error) {
	var req StartReceptionHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.PvzId == "" {
		return nil, ErrPvzIdRequired
	}
	return &req, nil
}
