package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
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

	err = h.merchantSvc.CreateMerchant(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// should return 201 if success
	w.WriteHeader(http.StatusCreated)
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

	customers, err := h.merchantSvc.GetMerchant(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// show response
	// fmt.Printf("GetMatch response: %+v\n", customers)

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = customers

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}