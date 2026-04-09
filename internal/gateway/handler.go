package gateway

import (
	"log"
	"net/http"

	"example.com/ecommerce/pkg/handler"
)

type gatewayService interface {
	Proxy(pr ProxyRequest) (http.HandlerFunc, error)
	HealthCheck() (*HealthResponse, error)
}

type Handler struct {
	handler.Base[gatewayService]
}

func NewHandler(service gatewayService) *Handler {
	return &Handler{
		Base: handler.Base[gatewayService]{Service: service},
	}
}

func (h *Handler) ProxyTo(target ServiceName) http.HandlerFunc {
	handlerFn, err := h.Service.Proxy(ProxyRequest{
		TargetService: target,
		StripPrefix:   "/api/v1",
		OnError:       h.BadGateway,
	})

	if err != nil {
		log.Printf("[gateway:handler] misconfigured proxy for %q: %v", target, err)
		return func(w http.ResponseWriter, r *http.Request) {
			h.BadGateway(w, "upstream service not configured")
		}
	}

	return handlerFn
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	result, err := h.Service.HealthCheck()

	if err != nil {
		h.ServiceUnavailable(w, err.Error())
		return
	}

	h.OK(w, result)
}
