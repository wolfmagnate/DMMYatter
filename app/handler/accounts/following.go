package accounts

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(usernameKey{}).(string)
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
	max_id, errMax := strconv.ParseInt(max_idStr, 10, 64)
	since_id, errSince := strconv.ParseInt(since_idStr, 10, 64)
	limit, errLimit := strconv.ParseInt(limitStr, 10, 64)

	if errMax != nil || errSince != nil || errLimit != nil || max_id < 0 || since_id < 0 || limit <= 0 {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	if max_idStr == "" {
		max_id = math.MaxInt64
	}

	if since_idStr == "" {
		since_id = math.MinInt64
	}

	if limitStr == "" {
		limit = math.MaxInt64
	}

	followees = followees.Filter(max_id, since_id, limit)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followees); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
