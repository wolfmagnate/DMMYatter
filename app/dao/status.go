package dao

import (
	"context"
	"time"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Status struct {
	ID        int64
	AccountID int64 `db:"account_id"`
	Content   string
	CreateAt  time.Time `db:"create_at"`
}

type status struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}
func (s *status) CreateStatus(ctx context.Context, status *object.Status) error {
	query := `INSERT INTO status (id, account_id, content, create_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, status.ID, status.PostedBy.ID, status.Content, status.CreateAt)
	return err
}

// Find a status with the specified id.
func (s *status) FindStatus(ctx context.Context, id int64) (*object.Status, error) {
	status := &Status{}
	retStatus := new(object.Status)

	// 投稿本体の情報を取得
	query := `SELECT * FROM status WHERE id = ?`
	err := s.db.QueryRowxContext(ctx, query, id).StructScan(status)
	if err != nil {
		return nil, err
	}
	retStatus.ID = status.ID
	retStatus.Content = status.Content
	retStatus.CreateAt = status.CreateAt

	// 投稿者を取得
	acc := NewAccount(s.db)
	postedby, err := acc.FindByID(ctx, status.AccountID)
	if err != nil {
		return nil, err
	}
	retStatus.PostedBy = postedby

	// attachmentテーブルのうち、status_idが等しいものに対応するmedia_idを全て取得する
	query = `SELECT media_id description FROM attachment WHERE status_id = ?`
	rows, err := s.db.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	attachedMedias := make([]*object.AttachedMedia, 0)
	med := NewMedia(s.db)
	for rows.Next() {
		var media_id int64
		var desc_text string
		err = rows.Scan(&media_id, &desc_text)
		if err != nil {
			return nil, err
		}
		media, err := med.FindMedia(ctx, media_id)
		if err != nil {
			return nil, err
		}
		attachedMedia := new(object.AttachedMedia)
		attachedMedia.Content = *media
		attachedMedia.Description = desc_text
		attachedMedias = append(attachedMedias, attachedMedia)
	}

	retStatus.AttachedMedias = attachedMedias
	return retStatus, nil
}

// Delete a status with the specified id.
func (s *status) DeleteStatus(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM attachment WHERE status_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM status WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
