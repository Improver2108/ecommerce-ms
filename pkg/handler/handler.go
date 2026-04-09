package handler

import (
	"encoding/json"
	"net/http"

	"example.com/ecommerce/pkg/response"
)

type Base[T any] struct {
	Service T
}

// Decode decodes the request body into v.
// Returns false and writes a 400 if decoding fails — handler should return immediately.
func (b *Base[T]) Decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		response.Err(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
func (b *Base[T]) OK(w http.ResponseWriter, data any) {
	response.Json(w, http.StatusOK, data)
}

func (b *Base[T]) Created(w http.ResponseWriter, data any) {
	response.Json(w, http.StatusCreated, data)
}

func (b *Base[T]) BadRequest(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusBadRequest, msg)
}

func (b *Base[T]) Unauthorized(w http.ResponseWriter, data any) {
	response.Json(w, http.StatusUnauthorized, data)
}

func (b *Base[T]) NotFound(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusNotFound, msg)
}

func (b *Base[T]) Conflict(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusConflict, msg)
}

func (b *Base[T]) InternalError(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusInternalServerError, msg)
}

func (b *Base[T]) BadGateway(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusBadGateway, msg)
}

func (b *Base[T]) ServiceUnavailable(w http.ResponseWriter, msg string) {
	response.Err(w, http.StatusServiceUnavailable, msg)
}
