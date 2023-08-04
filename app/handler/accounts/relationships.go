package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start relation")
	ctx := r.Context()
	self := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	usernames := r.URL.Query().Get("username")

	separatedNames := strings.Split(usernames, ",")

	specifiedAccounts := make([]*object.Account, 0, len(separatedNames))
	for _, name := range separatedNames {
		acc, err := h.ar.FindByUsername(ctx, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find account by username %s: %v", name, err), http.StatusInternalServerError)
			return
		}
		specifiedAccounts = append(specifiedAccounts, acc)
	}

	relationships, err := h.rr.GetRelationship(ctx, self, specifiedAccounts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get relationships: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode relationships to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
