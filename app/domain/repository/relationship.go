package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	// Follow a user. return true if succeeded
	FollowUser(ctx context.Context, follower *object.Account, followee *object.Account) error

	// Unfollow a user, return true if succeeded
	UnfollowUser(ctx context.Context, follower *object.Account, folowee *object.Account) error
}
