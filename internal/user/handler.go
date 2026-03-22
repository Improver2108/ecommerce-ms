package user

import (
	"net/http"

	"example.com/ecommerce/pkg/handler"
	"example.com/ecommerce/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	GetById(id string) (*User, error)
	UpdateUser(userId string, req UpdateUserRequest) (*User, error)
}

type Handler struct {
	handler.Base[userService]
}

func NewHandler(service userService) *Handler {
	return &Handler{
		Base: handler.Base[userService]{Service: service},
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

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || userID == "" {
		h.Unauthorized(w, "unauthorized")
	}

	var req UpdateUserRequest
	if !h.Decode(w, r, &req) {
		return
	}

	updated, err := h.Service.UpdateUser(userID, req)

	if err != nil {
		h.NotFound(w, err.Error())
	}

	h.OK(w, updated)
}
