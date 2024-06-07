package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/service"
	response "github.com/sprint-id/belimang/pkg/resp"
)

type itemHandler struct {
	itemSvc *service.ItemService
}

func newItemHandler(itemSvc *service.ItemService) *itemHandler {
	return &itemHandler{itemSvc}
}

func (h *itemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	merchantId := strings.Split(r.URL.Path, "/")[3]
	// fmt.Printf("id: %s\n", merchantId)
	var req dto.ReqAddItem
	var res dto.ResAddItem
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
	expectedFields := []string{"name", "price", "productCategory", "imageUrl"}
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

	res, err = h.itemSvc.AddItem(r.Context(), req, token.Subject(), merchantId)
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

func (h *itemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetItem

	param.ItemId = queryParams.Get("itemId")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Name = queryParams.Get("name")
	param.ProductCategory = queryParams.Get("productCategory")
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	records, err := h.itemSvc.GetItem(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessDataReponse{}
	successRes.Data = records
	successRes.Meta = response.Meta{
		Offset: param.Offset,
		Limit:  param.Limit,
		Total:  len(records),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status code to 201
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// The contains function checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
