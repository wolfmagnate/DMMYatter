package accounts

import (
	"encoding/json"
	"io"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var avatarData, headerData []byte
	for _, key := range []string{"avatar", "header"} {
		files, ok := r.MultipartForm.File[key]
		if !ok {
			continue
		}

		if len(files) != 1 {
			http.Error(w, "only one file is allowed", http.StatusBadRequest)
			return
		}

		file, err := files[0].Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		contentType := http.DetectContentType(data)
		if contentType != "image/jpeg" && contentType != "image/png" {
			http.Error(w, "file type not allowed", http.StatusBadRequest)
			return
		}

		if key == "avatar" {
			avatarData = data
		} else {
			headerData = data
		}
	}

	if noteValues, ok := r.MultipartForm.Value["note"]; ok && len(noteValues) > 0 {
		note := noteValues[0]
		account.Note = &note
	}

	displayNameValues, ok := r.MultipartForm.Value["display_name"]
	if !ok || len(displayNameValues) == 0 {
		http.Error(w, "display_name value not found", http.StatusBadRequest)
		return
	}
	displayName := displayNameValues[0]
	account.DisplayName = &displayName

	err = h.ar.UpdateAccountCredential(ctx, account, avatarData, headerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
