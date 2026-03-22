package auth

import (
	"net/http"

	"example.com/ecommerce/pkg/handler"
)

type AuthService interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
}

type Handler struct {
	handler.Base[AuthService]
}

func NewHandler(service AuthService) *Handler {
	return &Handler{
		Base: handler.Base[AuthService]{Service: service},
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if !h.Decode(w, r, &req) {
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		h.BadRequest(w, "name, email and password are required")
		return
	}

	resp, err := h.Service.Register(req)
	if err != nil {
		h.Conflict(w, err.Error())
		return
	}

	h.Created(w, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if !h.Decode(w, r, &req) {
		return
	}

	if req.Email == "" || req.Password == "" {
		h.BadRequest(w, "email and password are required")
		return
	}

	resp, err := h.Service.Login(req)
	if err != nil {
		h.Unauthorized(w, err.Error())
		return
	}

	h.OK(w, resp)
}
