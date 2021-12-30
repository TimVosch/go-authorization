package app

import (
	"context"
	"fmt"
	"goauthz/security"
)

var (
	UserPrincipalType = security.PrincipalType("user")
)

// User ...
type User struct {
	UserID int
	Rights []string
}

func (User) Type() security.PrincipalType {
	return UserPrincipalType
}

// UserAuthorizer ...
type UserAuthorizer struct {
}

func NewUserAuthorizer() *UserAuthorizer {
	return &UserAuthorizer{}
}

func (ua *UserAuthorizer) MustHaveRights(ctx context.Context, subject security.URN, rights ...string) error {
	user, ok := security.GetPrincipal(ctx).(*User)
	if !ok {
		fmt.Printf("UserAuthorizer does not work with %T\n", security.GetPrincipal(ctx))
		return security.ErrForbidden
	}

	if subject.Resource() == "user" {
		sid, err := security.URNToUserID(subject)
		if err != nil {
			fmt.Printf("Error converting URN to UserID: %v\n", err)
			return security.ErrForbidden
		}
		if sid == user.UserID {
			return nil
		}
	}

	return security.ErrForbidden
}
