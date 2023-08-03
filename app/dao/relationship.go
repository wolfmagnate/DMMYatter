package dao

import (
	"context"
	"database/sql"
	"errors"
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
	// INSERT INTO relationship (follower_id, followee_id) VALUES (3, 5);
	// すでに存在するかを判定
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, `select * from relationship where follower_id = ? and followee_id = ?`, follower.ID, followee.ID).StructScan(entity)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = r.db.ExecContext(ctx, `insert into relationship (follower_id, followee_id) values (?, ?)`, follower.ID, followee.ID)
		if err != nil {
			return err
		}
		return nil
	}

	// すでにフォロー関係が存在していた場合、適当にerrorを返す
	return nil // TODO
}
func (r *relationship) UnfollowUser(ctx context.Context, follower *object.Account, followee *object.Account) error {
	// DELETE FROM relationship WHERE follower_id = 3 AND followee_id = 5;
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, `SELECT * FROM relationship WHERE follower_id = ? AND followee_id = ?`, follower.ID, followee.ID).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("No such follow relationship exists")
		}
		return err
	}
	_, err = r.db.ExecContext(ctx, `DELETE FROM relationship WHERE follower_id = ? AND followee_id = ?`, follower.ID, followee.ID)
	if err != nil {
		return err
	}
	return nil
}
