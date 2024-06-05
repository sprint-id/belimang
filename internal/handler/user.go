package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/service"
)

type userHandler struct {
	userSvc *service.UserService
}

func newUserHandler(userSvc *service.UserService) *userHandler {
	return &userHandler{userSvc}
}

func (h *userHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAdminRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.RegisterAdmin(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqUserRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.RegisterUser(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.LoginAdmin(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.userSvc.LoginUser(r.Context(), req)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
