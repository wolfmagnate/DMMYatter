package dao

import (
	"context"
	"database/sql"
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
	query := `SELECT * FROM media WHERE id = ?`
	var existMedia object.Media
	err := m.db.QueryRowxContext(ctx, query, media.ID).StructScan(&existMedia)

	if err != sql.ErrNoRows {
		return err
	}

	query = `INSERT INTO media (id, type, url) VALUES (?, ?, ?)`
	_, err = m.db.ExecContext(ctx, query, media.ID, media.MediaType, media.URL)

	if err != nil {
		return err
	}
	return nil
}

// Find a media with the specified id
func (m *media) FindMedia(ctx context.Context, id int64) (*object.Media, error) {
	query := `SELECT * FROM media WHERE id = ?`
	var media object.Media
	err := m.db.QueryRowxContext(ctx, query, id).StructScan(&media)
	if err != nil {
		return nil, err
	}
	return &media, nil
}
