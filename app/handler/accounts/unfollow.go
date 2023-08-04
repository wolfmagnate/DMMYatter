package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(UsernameKey).(string)
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

	others := make([]*object.Account, 0)
	others = append(others, followee)
	relationships, err := h.rr.GetRelationship(ctx, follower, others)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships[0]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
