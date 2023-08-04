package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Media interface {
	// Create the specified media
	SaveMedia(ctx context.Context, media *object.Media) (int64, error)

	// Find a media with the specified id
	FindMedia(ctx context.Context, id int64) (*object.Media, error)
}
