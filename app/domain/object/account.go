package object

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	// The internal ID of the account
	ID int64 `json:"id,omitempty"`

	// The username of the account
	Username string `json:"username,omitempty"`

	// The username of the account
	PasswordHash string `json:"-" db:"password_hash"`

	// The account's display name
	DisplayName *string `json:"display_name,omitempty" db:"display_name"`

	// How many accounts follows the account
	FollowerCount int64 `json:"followers_count" db:"-"`

	// How many accounts the account follows
	FolloweeCount int64 `json:"following_count" db:"-"`

	// URL to the avatar image
	Avatar *string `json:"avatar,omitempty"`

	// URL to the header image
	Header *string `json:"header,omitempty"`

	// Biography of user
	Note *string `json:"note,omitempty"`

	// The time the account was created
	CreateAt time.Time `json:"create_at,omitempty" db:"create_at"`
}

type AccountGroup []*Account

// Check if given password is match to account's password
func (a *Account) CheckPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(pass)) == nil
}

// Hash password and set it to account object
func (a *Account) SetPassword(pass string) error {
	passwordHash, err := generatePasswordHash(pass)
	if err != nil {
		return fmt.Errorf("generate error: %w", err)
	}
	a.PasswordHash = passwordHash
	return nil
}

func generatePasswordHash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password failed: %w", err)
	}
	return string(hash), nil
}

func (a *Account) SetCreateAt() {
	a.CreateAt = time.Now()
}

func (accounts AccountGroup) Filter(max_id int64, since_id int64, limit int64) AccountGroup {
	var filteredAccounts AccountGroup
	count := int64(0)
	for _, account := range accounts {
		if account.ID > since_id && account.ID < max_id {
			filteredAccounts = append(filteredAccounts, account)
			count++
			if count >= limit {
				break
			}
		}
	}
	return filteredAccounts
}
