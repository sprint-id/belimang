package handler

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sprint-id/belimang/internal/cfg"
	"github.com/sprint-id/belimang/internal/service"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	merchantH := newMerchantHandler(h.service.Merchant)
	itemH := newItemHandler(h.service.Item)
	estimateH := newEstimateHandler(h.service.Estimate)
	fileH := newFileHandler(h.cfg)

	r.Use(middleware.RedirectSlashes)

	r.Post("/admin/register", userH.RegisterAdmin)
	r.Post("/admin/login", userH.LoginAdmin)
	r.Post("/users/register", userH.RegisterUser)
	r.Post("/users/login", userH.LoginUser)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/admin/merchants", merchantH.CreateMerchant)
		r.Get("/admin/merchants", merchantH.GetMerchant)
		r.Get("/merchants/nearby/{lat},{long}", merchantH.GetNearbyMerchant)
		r.Post("/admin/merchants/{merchantId}/items", itemH.AddItem)
		r.Get("/admin/merchants/{merchantId}/items", itemH.GetItem)

		r.Post("/users/estimate", estimateH.CreateEstimate)

		r.Post("/image", fileH.Upload)
	})
}
