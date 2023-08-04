package accounts

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"yatter-backend-go/app/handler/util"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(UsernameKey).(string)
	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followees, err := h.ar.FindFolloweeOfAccount(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	max_idStr, since_idStr, limitStr := r.URL.Query().Get("max_id"), r.URL.Query().Get("since_id"), r.URL.Query().Get("limit")
	max_id, err := util.ParseInt64QueryParam(max_idStr, math.MaxInt64, func(v int64) bool { return v >= 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	since_id, err := util.ParseInt64QueryParam(since_idStr, 0, func(v int64) bool { return v >= 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := util.ParseInt64QueryParam(limitStr, 40, func(v int64) bool { return v > 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	followees = followees.Filter(max_id, since_id, limit)
	fmt.Println(followees)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followees); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
