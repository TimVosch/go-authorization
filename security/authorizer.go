package security

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type PrincipalType string

// Authorizer ...
type Authorizer interface {
	MustHaveRights(ctx context.Context, subject URN, rights ...string) error
}

type AuthorizationRegistry struct {
	authorizers map[PrincipalType]Authorizer
}

func NewAuthorizerRegistry() *AuthorizationRegistry {
	return &AuthorizationRegistry{
		authorizers: make(map[PrincipalType]Authorizer),
	}
}

func (reg *AuthorizationRegistry) Add(principle PrincipalType, authorizer Authorizer) *AuthorizationRegistry {
	reg.authorizers[principle] = authorizer
	return reg
}

func (reg *AuthorizationRegistry) Get(principle PrincipalType) (Authorizer, bool) {
	authorizer, ok := reg.authorizers[principle]
	return authorizer, ok
}

func (reg *AuthorizationRegistry) MustHaveRights(ctx context.Context, subject URN, rights ...string) error {
	principal := GetPrincipal(ctx)
	if principal == nil {
		fmt.Println("Unauthenticated")
		return ErrUnauthenticated
	}
	mech, ok := reg.authorizers[principal.Type()]
	if !ok {
		fmt.Println("Authenticated but no authorizer found for the principal type")
		return ErrForbidden
	}
	return mech.MustHaveRights(ctx, subject, rights...)
}
