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

type merchantHandler struct {
	merchantSvc *service.MerchantService
}

func newMerchantHandler(merchantSvc *service.MerchantService) *merchantHandler {
	return &merchantHandler{merchantSvc}
}

// {
// 	"name": "", // not null | minLength 2 | maxLength 30
// 	"merchantCategory": "", /** enum of:
// 	- `SmallRestaurant`
// 	- `MediumRestaurant`
// 	- `LargeRestaurant`
// 	- `MerchandiseRestaurant`
// 	- `BoothKiosk`
// 	- `ConvenienceStore`
// 		*/
// 	"imageUrl": "", // not null | should be image url
//   "location": {
//     "lat": 1, // not null | float
//     "long": 1  // not null | float
//   }
// }

func (h *merchantHandler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateMerchant
	var res dto.ResCreateMerchant
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		fmt.Printf("error Decode: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"name", "merchantCategory", "imageUrl", "location"}
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

	// show request
	fmt.Printf("CreateMerchant request: %+v\n", req)

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.merchantSvc.CreateMerchant(r.Context(), req, token.Subject())
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

func (h *merchantHandler) GetMerchant(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetMerchant

	param.MerchantId = queryParams.Get("merchantId")
	param.Name = queryParams.Get("name")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.MerchantCategory = queryParams.Get("merchantCategory")
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	merchants, err := h.merchantSvc.GetMerchant(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessPageReponse{}
	successRes.Message = "success"
	successRes.Data = merchants
	successRes.Meta = response.Meta{
		Offset: param.Offset,
		Limit:  param.Limit,
		Total:  len(merchants),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status code to 201
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *merchantHandler) GetNearbyMerchant(w http.ResponseWriter, r *http.Request) {
	// get lat and long from URL
	var latFloat, longFloat float64
	latAndLong := strings.Split(r.URL.Path, "/")[3]
	fmt.Printf("latAndLong: %s\n", latAndLong)
	lat := strings.Split(latAndLong, ",")[0]
	long := strings.Split(latAndLong, ",")[1]
	fmt.Printf("lat: %s, long: %s\n", lat, long)

	// check if lat long not a number
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "lat is not a number", http.StatusBadRequest)
		return
	}

	longFloat, err = strconv.ParseFloat(long, 64)
	if err != nil {
		http.Error(w, "long is not a number", http.StatusBadRequest)
		return
	}

	// show lat and long in float
	fmt.Printf("after parse lat: %.9f, long: %.9f\n", latFloat, longFloat)

	// Query params
	queryParams := r.URL.Query()
	var param dto.ParamGetNearbyMerchant

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

	merchants, err := h.merchantSvc.GetNearbyMerchant(r.Context(), param, token.Subject(), latFloat, longFloat)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessDataReponse{}
	successRes.Data = merchants
	successRes.Meta = response.Meta{
		Offset: param.Offset,
		Limit:  param.Limit,
		Total:  len(merchants),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status code to 201
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
