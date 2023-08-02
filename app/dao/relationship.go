package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Relationship struct {
	FollowerID int64 `db:"follower_id"`
	FolloweeID int64 `db:"followee_id"`
}

type relationship struct {
	db *sqlx.DB
}

func NewRelationship(db *sqlx.DB) repository.Relationship {
	return &relationship{db: db}
}

func (r *relationship) FollowUser(ctx context.Context, follower *object.Account, followee *object.Account) error {

}

func (r *relationship) UnfollowUser(ctx context.Context, follower *object.Account, folowee *object.Account) error {

}
