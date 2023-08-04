package object

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

type Media struct {
	// internal id
	ID int64 `json:"media_id" db:"id"`

	// The URL where the content of this media is saved.
	URL string `json:"media_url" db:"url"`

	// Type of this media. For example image, video, sound.
	MediaType string `json:"-" db:"type"`
}

func NewDummyHostMedia(data []byte) *Media {
	m := Media{}

	hasher := sha256.New()
	hasher.Write(data)
	dummyURL := hex.EncodeToString(hasher.Sum(nil))
	m.URL = fmt.Sprintf("https://dummyimage.com/%s", dummyURL)

	m.MediaType = http.DetectContentType(data)
	return &m
}
