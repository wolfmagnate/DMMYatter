package statuses

import (
	"context"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	sr repository.Status
	mr repository.Media
}

// Create Handler for `/v1/statuses/`
func NewRouter(ar repository.Account, mr repository.Media, sr repository.Status) http.Handler {
	r := chi.NewRouter()

	h := &handler{sr, mr}
	r.With(auth.Middleware(ar)).Post("/", h.Post)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(idContext)
		r.Get("/", h.Get)
		r.With(auth.Middleware(ar)).Delete("/", h.Delete)
	})

	return r
}

type idKey struct{}

var IDKey = idKey{}

func idContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), IDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
