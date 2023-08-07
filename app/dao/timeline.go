package dao

import (
	"context"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type timeline struct {
	db *sqlx.DB
}

func NewTimeline(db *sqlx.DB) repository.Timeline {
	return &timeline{db: db}
}

func (t *timeline) GetHome(ctx context.Context, id int64, maxID int64, sinceID int64, limit int64) ([]*object.Status, error) {
	var statusesID []int64

	query := `
		SELECT s.id 
		FROM status AS s
		JOIN (
			SELECT followee_id 
			FROM relationship 
			WHERE follower_id = ?
		) AS followings ON followings.followee_id = s.account_id
		WHERE s.id > ? AND s.id < ?
		LIMIT ?
	`

	rows, err := t.db.QueryxContext(ctx, query, id, sinceID, maxID, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		statusesID = append(statusesID, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	s := NewStatus(t.db)
	timelineStatuses := make([]*object.Status, 0)
	for _, id := range statusesID {
		item, err := s.FindStatus(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error creating timeline: %v", err)
		}
		timelineStatuses = append(timelineStatuses, item)
	}

	return timelineStatuses, nil
}

func (t *timeline) GetPublic(ctx context.Context, maxID int64, sinceID int64, limit int64) ([]*object.Status, error) {
	var statusesID []int64

	query := `
		SELECT s.id 
		FROM status AS s
		WHERE s.id > ? AND s.id < ?
		LIMIT ?
	`

	rows, err := t.db.QueryxContext(ctx, query, sinceID, maxID, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		statusesID = append(statusesID, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	s := NewStatus(t.db)
	timelineStatuses := make([]*object.Status, 0)
	for _, id := range statusesID {
		item, err := s.FindStatus(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error creating timeline: %v", err)
		}
		timelineStatuses = append(timelineStatuses, item)
	}

	return timelineStatuses, nil
}
