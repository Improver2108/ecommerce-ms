package user

import (
	"net/http"

	"example.com/ecommerce/pkg/handler"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	handler.Base[Service]
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Base: handler.Base[Service]{Service: service},
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		h.BadRequest(w, "missing user id")
		return
	}

	user, err := h.Service.GetById(id)

	if err != nil {
		h.NotFound(w, err.Error())
		return
	}

	h.OK(w, user)
}
