package accounts

import (
	"net/http"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(usernameKey{}).(string)
	authUsername := ctx.Value(auth.AuthUsernameKey{}).(string)

	// クエリパラメータからfolloweeを取得する
	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 認証ヘッダを元にfollowerを取得する
	follower, err := h.ar.FindByUsername(ctx, authUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.rr.FollowUser(ctx, follower, followee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

}
