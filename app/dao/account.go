package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}

	return entity, nil
}

func (r *account) FindByID(ctx context.Context, id int64) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where id = ?", id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}

	return entity, nil
}

// Find follower of specified account
func (r *account) FindFollowerOfAccount(ctx context.Context, followee *object.Account) (object.AccountGroup, error) {
	rows, err := r.db.QueryxContext(ctx, "select acc.* from account as acc join (select * from relationship where followee_id = ?) as rel on acc.id = rel.follower_id", followee.ID)
	if err != nil {
		return nil, err
	}
	followers := make([]*object.Account, 0)
	for rows.Next() {
		entity := new(object.Account)
		err = rows.StructScan(entity)
		if err != nil {
			return nil, err
		}
		followers = append(followers, entity)
	}
	return followers, nil
}

// Find followee of specified account
func (r *account) FindFolloweeOfAccount(ctx context.Context, follower *object.Account) (object.AccountGroup, error) {
	rows, err := r.db.QueryxContext(ctx, "select acc.* from account as acc join (select * from relationship where follower_id = ?) as rel on acc.id = rel.followeee_id", follower.ID)
	if err != nil {
		return nil, err
	}
	followees := make([]*object.Account, 0)
	for rows.Next() {
		entity := new(object.Account)
		err = rows.StructScan(entity)
		if err != nil {
			return nil, err
		}
		followees = append(followees, entity)
	}
	return followees, nil
}

// Update account information. Return true if succeeded
func (r *account) UpdateAccountCredential(ctx context.Context, account *object.Account) error {
	_, err := r.db.ExecContext(ctx, "update account set display_name = ?, note = ?, avatar = ?, header = ? where id = ?", account.DisplayName, account.Note, account.Avatar, account.Header, account.ID)
	if err != nil {
		return err
	}
	return nil
}

// Create a new account. Return the created account
func (r *account) CreateNewAccount(ctx context.Context, account *object.Account) (*object.Account, error) {
	_, err := r.db.ExecContext(ctx, "insert into account (username, password_hash, display_name, avatar, header, note, create_at) values (?, ?, ?, ?, ?, ?, ?, ?)", account.Username, account.PasswordHash, account.DisplayName, account.Avatar, account.Header, account.Note, account.CreateAt)
	if err != nil {
		return nil, err
	}
	entity := new(object.Account)
	err = r.db.QueryRowxContext(ctx, "select * from account where id = ?", account.ID).StructScan(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
