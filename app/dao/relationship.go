package dao

import (
	"context"
	"errors"
	"fmt"
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
	var count int

	err := r.db.GetContext(ctx, &count, "SELECT count(*) FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID)
	if err != nil {
		return fmt.Errorf("failed to query relationship from db: %w", err)
	}

	if count == 0 {
		tx, err := r.db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, `INSERT INTO relationship (follower_id, followee_id) values (?, ?)`, follower.ID, followee.ID)
		if err != nil {
			return fmt.Errorf("failed to insert into relationship: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to add follow because already followed: %w", errors.New("already followed"))
}

func (r *relationship) UnfollowUser(ctx context.Context, follower *object.Account, followee *object.Account) error {
	var count int

	err := r.db.GetContext(ctx, &count, "SELECT count(*) FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID)
	if err != nil {
		return fmt.Errorf("failed to query relationship from db: %w", err)
	}

	if count != 0 {
		tx, err := r.db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, `DELETE FROM relationship WHERE follower_id = ? AND followee_id = ?`, follower.ID, followee.ID)
		if err != nil {
			return fmt.Errorf("failed to delete from relationship: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		return fmt.Errorf("failed to unfollow because follow relationship does not exist: %w", errors.New("follow relationship does not exist"))
	}

	return nil
}
