package media

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	mr repository.Media
}

// Create Handler for `/v1/media/`
func NewRouter(mr repository.Media) http.Handler {
	r := chi.NewRouter()

	h := &handler{mr}
	r.Post("/", h.Post)
	return r
}
