package statuses

import (
	"net/http"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value(IDKey).(int64)

	err := h.sr.DeleteStatus(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
