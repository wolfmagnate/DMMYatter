package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type media struct {
	db *sqlx.DB
}

func NewMedia(db *sqlx.DB) repository.Media {
	return &media{db: db}
}

// Create the specified media
func (m *media) SaveMedia(ctx context.Context, media *object.Media) error {
	// insert into media (id, type, url) values (1, "image", "https://example.com/image1.jpg")
}

// Find a media with the specified id
func (m *media) FindMedia(ctx context.Context, id int64) (*object.Media, error) {
	// select * from media where id = 引数のid
}
