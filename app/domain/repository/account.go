package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified id
	FindByID(ctx context.Context, id int64) (*object.Account, error)
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	// Find follower of specified account
	FindFollowerOfAccount(ctx context.Context, followee *object.Account) (object.AccountGroup, error)
	// Find followee of specified account
	FindFolloweeOfAccount(ctx context.Context, follower *object.Account) (object.AccountGroup, error)

	// Update account information.
	UpdateAccountCredential(ctx context.Context, account *object.Account, avatarData []byte, headerData []byte) error
	// Create a new account.
	CreateNewAccount(ctx context.Context, account *object.Account) (*object.Account, error)
}
