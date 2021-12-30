package usersession

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"goauthz/security"

	"github.com/jmoiron/sqlx"
)

var (
	Type = security.AuthType("basic")

	_ security.Authenticator = (*SessionManager)(nil)
)

type PrincipalCreator func(context.Context, *Session) (security.Principal, error)

// SessionManager ...
type SessionManager struct {
	db               *sqlx.DB
	principalCreator PrincipalCreator
}

func New(db *sqlx.DB, principalCreator PrincipalCreator) *SessionManager {
	return &SessionManager{
		db:               db,
		principalCreator: principalCreator,
	}
}

// Authenticate implements security.Authenticator interface
func (s *SessionManager) Authenticate(ctx context.Context, cred security.Credential) (context.Context, error) {
	if cred.Type != Type {
		return ctx, errors.New("invalid credential type")
	}

	session, err := getSession(s.db, cred.Secret)
	if err != nil {
		return ctx, err
	}

	principal, err := s.principalCreator(ctx, session)
	if err != nil {
		return ctx, err
	}

	return security.WithPrincipal(ctx, principal), nil
}

func (s *SessionManager) CreateSession(data SessionData) (*Session, security.Credential, error) {
	var c security.Credential
	c.Type = Type

	key := randomKey()
	if key == "" {
		return nil, c, errors.New("failed to generate key")
	}

	sess, err := createSession(s.db, key, data)
	if err != nil {
		return nil, c, err
	}

	c.Secret = key
	return sess, c, nil
}

func randomKey() string {
	var data = make([]byte, 16)
	_, err := rand.Read(data)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(data)
}
