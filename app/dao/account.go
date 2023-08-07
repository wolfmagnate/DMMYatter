package dao

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

/*
	現在のdaoは、ハンドラで使われることを想定して、かなりリッチな（単なるDB:オブジェクト変換以上の）仕事をしている
	そもそもDBへの入出力とサーバーとの入出力はそれぞれ独立した別の場所への入出力であるので、同じ場所で行うべきではない
	DAOをリッチにして戻り値を完成済み状態に近づけないならば、結果を受け取るハンドラ側の処理が増えて、結局fat controller状態になりドメインロジックが漏れる
	結局これは間にドメインを挟まずに直接DBとハンドラを接続したことによる弊害である。
	間に必ずドメインがあると仮定できるならば、DAOでは単にDBからドメインのオブジェクトのベースとなる情報を返すことに専念できる
*/

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
func (r *account) UpdateAccountCredential(ctx context.Context, account *object.Account, avatarData []byte, headerData []byte) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	if account.Avatar != nil {
		avatarFileName, err := r.saveImage(avatarData, "./postedImages")
		if err != nil {
			return fmt.Errorf("failed to save avatar image: %w", err)
		}
		account.Avatar = &avatarFileName
	}

	if account.Header != nil {
		headerFileName, err := r.saveImage(headerData, "./postedImages")
		if err != nil {
			return fmt.Errorf("failed to save header image: %w", err)
		}

		account.Header = &headerFileName
	}

	_, err = tx.ExecContext(ctx, "update account set display_name = ?, note = ?, avatar = ?, header = ? where id = ?", account.DisplayName, account.Note, account.Avatar, account.Header, account.ID)
	if err != nil {
		return fmt.Errorf("failed to update account db: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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

func (r *account) saveImage(imageData []byte, dir string) (string, error) {
	hash := sha256.Sum256(imageData)

	contentType := http.DetectContentType(imageData)

	var extension string
	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	fileName := fmt.Sprintf("%x%s", hash, extension)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}

	return fileName, nil
}
