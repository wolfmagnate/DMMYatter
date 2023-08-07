package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yatter-backend-go/app/domain/object"
)

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("start media post")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mediaFile, ok := r.MultipartForm.File["file"]
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(mediaFile) != 1 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, fileHeader := range mediaFile {
		f, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// bを適当なところにホストする処理が必要だが、やらない
		// 適当なURLにホストしたことにする
		newMedia := object.NewDummyHostMedia(b)
		id, err := h.mr.SaveMedia(ctx, newMedia)
		newMedia.ID = id

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(newMedia); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
