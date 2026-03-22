package authapi

import (
	"example.com/ecommerce/internal/auth"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *auth.Handler) {
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
}
