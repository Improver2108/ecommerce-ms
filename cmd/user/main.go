package main

import (
	"fmt"
	"net/http"

	userapi "example.com/ecommerce/api/user"
	"example.com/ecommerce/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	svc := user.NewService()
	handler := user.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userapi.RegisterRoutes(r, handler)
	})

	fmt.Println("User service running on :8082")
	http.ListenAndServe(":8082", r)
}
