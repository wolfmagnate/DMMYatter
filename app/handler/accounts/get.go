package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value(UsernameKey).(string)
	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(account)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
