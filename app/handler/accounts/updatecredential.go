package accounts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

func (h *handler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	// クエリの情報に従ってaccountの内容を更新する
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ar.UpdateAccountCredential(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "# MultipartForm.Value")
	for k, v := range r.MultipartForm.Value {
		fmt.Fprintln(w, k, v)
	}

	fmt.Fprintln(w, "# MultipartForm.File")
	for k, v := range r.MultipartForm.File {
		fmt.Fprintln(w, "## Key:", k)
		for _, fh := range v {
			fmt.Fprintln(w, "## Filename:", fh.Filename)
			f, err := fh.Open()
			if err != nil {
				fmt.Fprintln(w, err)
				break
			}
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				fmt.Fprintln(w, err)
				break
			}
			fmt.Fprintln(w, "## Body:")
			fmt.Fprintln(w, string(b))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
