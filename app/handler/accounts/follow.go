package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(UsernameKey).(string)
	follower := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.rr.FollowUser(ctx, follower, followee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	others := make([]*object.Account, 0)
	others = append(others, followee)

	// この辺りが本来別々であるはずのハンドラの入出力とDBの入出力がドメインを介さず直接行われているために、
	// ある程度ちゃんと動作するモックDBがないとロジックの検証ができないから辛そう
	// インターフェースによって具体的なDB製品や実装方式とかへの依存は外せているけれど
	// 具体的（DBに近い）インターフェースに依存してること自体は変わっていない
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
