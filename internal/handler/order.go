package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/service"
)

type orderHandler struct {
	orderSvc *service.OrderService
}

func newOrderHandler(orderSvc *service.OrderService) *orderHandler {
	return &orderHandler{orderSvc}
}

func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateOrder
	var res dto.ResCreateOrder
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		fmt.Printf("error Decode: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the payload is empty
	if len(jsonData) == 0 {
		http.Error(w, "empty payload", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"calculatedEstimateId"}
	for key := range jsonData {
		if !contains(expectedFields, key) {
			http.Error(w, "unexpected field in request body: "+key, http.StatusBadRequest)
			return
		}
	}

	// Convert the jsonData map into the req struct
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Printf("error json Marshal: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		fmt.Printf("error json Unmarshal: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.orderSvc.CreateOrder(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set HTTP status code to 201
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *orderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	var res []dto.ResGetOrderHistory
	queryParams := r.URL.Query()
	var param dto.ParamGetOrderHistory

	// Decode request query into the param struct
	param.MerchantId = queryParams.Get("merchantId")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Name = queryParams.Get("name")
	param.MerchantCategory = queryParams.Get("merchantCategory")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.orderSvc.GetOrderHistory(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status code to 200
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
