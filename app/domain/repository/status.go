package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Create the specified status.
	CreateStatus(ctx context.Context, status *object.Status) error
	// Find a status with the specified id.
	FindStatus(ctx context.Context, id int64) (*object.Status, error)
	// Delete a status with the specified id.
	DeleteStatus(ctx context.Context, id int64) (*object.Status, error)
}
