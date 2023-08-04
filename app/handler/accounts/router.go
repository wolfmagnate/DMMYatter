package accounts

import (
	"context"
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	ar repository.Account
	rr repository.Relationship
}

// Create Handler for `/v1/accounts/`
func NewRouter(ar repository.Account, rr repository.Relationship) http.Handler {

	r := chi.NewRouter()

	h := &handler{ar, rr}
	r.Post("/", h.Create)
	r.With(auth.Middleware(ar)).Post("/update_credentials", h.UpdateCredential)
	r.With(auth.Middleware(ar)).Get("/relationships", h.Relationships)
	r.Route("/{username}", func(r chi.Router) {
		r.Use(usernameContext)
		r.Get("/", h.Get)
		r.With(auth.Middleware(ar)).Post("/follow", h.Follow)
		r.With(auth.Middleware(ar)).Post("/unfollow", h.Unfollow)
		r.Get("/following", h.Following)
		r.Get("/followers", h.Followers)
	})

	return r
}

type usernameKey struct{}

var UsernameKey = usernameKey{}

func usernameContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
