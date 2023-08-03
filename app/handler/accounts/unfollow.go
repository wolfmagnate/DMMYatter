package accounts

import (
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(usernameKey{}).(string)
	follower := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	// クエリパラメータからfolloweeを取得する
	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.rr.UnfollowUser(ctx, follower, followee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

}
