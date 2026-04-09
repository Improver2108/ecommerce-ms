package main

import (
	"fmt"
	"net/http"

	authapi "example.com/ecommerce/api/auth"
	"example.com/ecommerce/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	svc := auth.NewService()
	handler := auth.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authapi.RegisterRoutes(r, handler)

	fmt.Println("Auth service running on :8081")
	http.ListenAndServe(":8081", r)
}
