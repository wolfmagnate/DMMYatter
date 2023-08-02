package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type status struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// Create the specified status.
func (*status) CreateStatus(ctx context.Context, status *object.Status) error {
	// insert into status
	// attachmentがなければ生成
}

// Find a status with the specified id.
func (*status) FindStatus(ctx context.Context, id int64) (*object.Status, error) {
	// statusを取り出し、ついでに対応するattachment_bindingを取り出し、さらにはattachmentを結合する
}

// Delete a status with the specified id.
func (*status) DeleteStatus(ctx context.Context, id int64) (*object.Status, error) {
	// statusを削除し、ついでに対応するattachment_bindingを削除する
}
