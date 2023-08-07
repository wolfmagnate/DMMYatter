package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

type PostRequest struct {
	Status string
	Medias []PostAttachment
}

type PostAttachment struct {
	MediaID     int64 `json:"media_id"`
	Description string
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postAccount := ctx.Value(auth.AuthUsernameKey).(*object.Account)

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newStatus := new(object.Status)
	newStatus.Content = req.Status
	newStatus.PostedBy = postAccount
	newStatus.SetCreateAt()
	newStatus.AttachedMedias = make([]*object.AttachedMedia, 0)
	for _, attached := range req.Medias {
		attachedMedia, err := h.mr.FindMedia(ctx, attached.MediaID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newStatus.AttachedMedias = append(newStatus.AttachedMedias, &object.AttachedMedia{Content: *attachedMedia, Description: attached.Description})
	}

	h.sr.CreateStatus(ctx, newStatus)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(newStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
