package userapi

import (
	"example.com/ecommerce/internal/user"
	"example.com/ecommerce/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *user.Handler) {
	// user — protected (JWT middleware applied at gateway level)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/users/{id}", h.GetUser)
	})
}
