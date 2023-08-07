package timelines

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	mr repository.Media
	sr repository.Status
	tr repository.Timeline
}

// Create Handler for `/v1/timelines/`
func NewRouter(ar repository.Account, mr repository.Media, sr repository.Status, tr repository.Timeline) http.Handler {
	r := chi.NewRouter()

	h := &handler{mr, sr, tr}
	r.With(auth.Middleware(ar)).Get("/home", h.Home)
	r.Get("/public", h.Public)

	return r
}
