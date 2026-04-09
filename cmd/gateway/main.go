package main

import (
	"fmt"
	"net/http"

	gatewayapi "example.com/ecommerce/api/gateway"
	"example.com/ecommerce/internal/gateway"

	"github.com/go-chi/chi/v5"
)

func main() {
	svc := gateway.NewGatewayService(gateway.DefaultConfigs())
	handler := gateway.NewHandler(svc)

	r := chi.NewRouter()
	gatewayapi.RegisterRoutes(r, handler)

	fmt.Println("Gateway service running on :8080")
	http.ListenAndServe(":8080", r)

}
