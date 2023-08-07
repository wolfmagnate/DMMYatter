package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Timeline interface {
	GetHome(ctx context.Context, id int64, maxID int64, sinceID int64, limit int64) ([]*object.Status, error)
	GetPublic(ctx context.Context, maxID int64, sinceID int64, limit int64) ([]*object.Status, error)
}
