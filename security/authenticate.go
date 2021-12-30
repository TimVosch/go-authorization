package security

import (
	"context"
	"errors"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

type AuthType string

// Authenticator ...
type Authenticator interface {
	Authenticate(context.Context, Credential) (context.Context, error)
}

// AuthenticationRegistry ...
type AuthenticationRegistry struct {
	authenticators map[AuthType]Authenticator
}

func NewAuthenticationRegistery() *AuthenticationRegistry {
	return &AuthenticationRegistry{
		authenticators: make(map[AuthType]Authenticator),
	}
}

func (reg *AuthenticationRegistry) Add(typ AuthType, mech Authenticator) *AuthenticationRegistry {
	reg.authenticators[typ] = mech
	return reg
}

func (reg *AuthenticationRegistry) Get(typ AuthType) (Authenticator, bool) {
	mech, ok := reg.authenticators[typ]
	return mech, ok
}

func (reg *AuthenticationRegistry) Authenticate(ctx context.Context, c Credential) (context.Context, error) {
	mech, ok := reg.authenticators[c.Type]
	if !ok {
		return ctx, errors.New("no mechanism found")
	}

	return mech.Authenticate(ctx, c)
}
