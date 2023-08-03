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
func (r *account) countFollowersAndFollowees(ctx context.Context, entity *object.Account) error {
	err := r.db.QueryRowxContext(ctx, "select count(*) from relationship where followee_id = ?", entity.ID).Scan(&entity.FollowerCount)
	if err != nil {
		return err
	}
	err = r.db.QueryRowxContext(ctx, "select count(*) from relationship where follower_id = ?", entity.ID).Scan(&entity.FolloweeCount)
	if err != nil {
		return err
	}
	return nil
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

	err = r.countFollowersAndFollowees(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to count followers and followees: %w", err)
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

	err = r.countFollowersAndFollowees(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to count followers and followees: %w", err)
	}

	return entity, nil
}

// Find follower of specified account
func (r *account) FindFollowerOfAccount(ctx context.Context, followee *object.Account) (object.AccountGroup, error) {
	rows, err := r.db.QueryxContext(ctx, "select acc.* from account as acc join (select * from relationship where followee_id = ?) as rel on acc.id = rel.follower_id", followee.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query followers from db: %w", err)
	}

	followers := make([]*object.Account, 0)
	for rows.Next() {
		entity := new(object.Account)
		err = rows.StructScan(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan followers from db: %w", err)
		}
		followers = append(followers, entity)
	}

	return followers, nil
}

// Find followee of specified account
func (r *account) FindFolloweeOfAccount(ctx context.Context, follower *object.Account) (object.AccountGroup, error) {
	rows, err := r.db.QueryxContext(ctx, "select acc.* from account as acc join (select * from relationship where follower_id = ?) as rel on acc.id = rel.followee_id", follower.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query followees from db: %w", err)
	}

	followees := make([]*object.Account, 0)
	for rows.Next() {
		entity := new(object.Account)
		err = rows.StructScan(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan followees from db: %w", err)
		}
		followees = append(followees, entity)
	}

	return followees, nil
}

// Update account information. Return true if succeeded
func (r *account) UpdateAccountCredential(ctx context.Context, account *object.Account) error {
	_, err := r.db.ExecContext(ctx, "update account set display_name = ?, note = ?, avatar = ?, header = ? where id = ?", account.DisplayName, account.Note, account.Avatar, account.Header, account.ID)
	if err != nil {
		return fmt.Errorf("failed to update account db: %w", err)
	}
	return nil
}

// Create a new account. Return the created account
func (r *account) CreateNewAccount(ctx context.Context, account *object.Account) (*object.Account, error) {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "insert into account (username, password_hash, display_name, avatar, header, note, create_at) values (?, ?, ?, ?, ?, ?, ?)", account.Username, account.PasswordHash, account.DisplayName, account.Avatar, account.Header, account.Note, account.CreateAt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tx.Commit()
	entity := new(object.Account)
	err = r.db.QueryRowxContext(ctx, "select * from account where username = ?", account.Username).StructScan(entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return entity, nil
}
