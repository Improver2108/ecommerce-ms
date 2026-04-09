package gatewayapi

import (
	"example.com/ecommerce/internal/gateway"
	md "example.com/ecommerce/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(r chi.Router, h *gateway.Handler) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/health", h.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.ProxyTo(gateway.ServiceAuth))
			r.Post("/login", h.ProxyTo(gateway.ServiceAuth))
		})

		r.Group(func(r chi.Router) {
			r.Use(md.JWTMiddleware)
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", h.ProxyTo(gateway.ServiceUser))
			})
		})
	})
}
