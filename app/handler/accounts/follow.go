package accounts

import (
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(usernameKey{}).(string)
	follower := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	fmt.Println("find followee")
	// クエリパラメータからfolloweeを取得する
	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("start follow")
	err = h.rr.FollowUser(ctx, follower, followee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

}
