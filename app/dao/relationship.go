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

func (r *relationship) GetRelationship(ctx context.Context, self *object.Account, others []*object.Account) ([]*object.Relationship, error) {
	results := make([]*object.Relationship, 0)
	for _, other := range others {
		result := new(object.Relationship)
		result.OtherID = other.ID

		var followedByCount int
		err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE followee_id = ? AND follower_id = ?", self.ID, other.ID).Scan(&followedByCount)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query for checking if user with ID %d is followed by user with ID %d: %w", self.ID, other.ID, err)
		}
		result.FollowedBy = (followedByCount == 1)

		var followingCount int
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE followee_id = ? AND follower_id = ?", other.ID, self.ID).Scan(&followingCount)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query for checking if user with ID %d is following user with ID %d: %w", self.ID, other.ID, err)
		}
		result.Following = (followingCount == 1)

		results = append(results, result)
	}

	return results, nil
}
