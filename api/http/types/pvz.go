package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type OpenPvzHandlerRequest struct {
	City string `json:"city"`
}

func CreateOpenPvzHandlerRequest(r *http.Request) (*OpenPvzHandlerRequest, error) {
	var req OpenPvzHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.City != "Москва" && req.City != "Санкт-Петербург" && req.City != "Казань" {
		return nil, ErrInvalidCity
	}
	return &req, nil
}

type ListPvzHandlerRequest struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

func CreateListPvzHandlerRequest(r *http.Request) (*ListPvzHandlerRequest, error) {
	var req ListPvzHandlerRequest
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if startDateStr != "" {
		startTime, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid startDate format, expected RFC3339")
		}
		req.StartDate = &startTime
	}
	if endDateStr != "" {
		endTime, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format, expected RFC3339")
		}
		req.EndDate = &endTime
	}
	req.Page = 1
	req.Limit = 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			req.Page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = l
		}
	}
	return &req, nil
}
