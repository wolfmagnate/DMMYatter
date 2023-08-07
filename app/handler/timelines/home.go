package timelines

import (
	"encoding/json"
	"math"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/util"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postAccount := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	max_idStr, since_idStr, limitStr := r.URL.Query().Get("max_id"), r.URL.Query().Get("since_id"), r.URL.Query().Get("limit")
	maxID, err := util.ParseInt64QueryParam(max_idStr, math.MaxInt64, func(v int64) bool { return v >= 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sinceID, err := util.ParseInt64QueryParam(since_idStr, 0, func(v int64) bool { return v >= 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := util.ParseInt64QueryParam(limitStr, 40, func(v int64) bool { return v > 0 })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultTimeline, err := h.tr.GetHome(ctx, postAccount.ID, maxID, sinceID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resultTimeline); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
