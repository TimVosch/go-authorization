package usersession

import (
	"database/sql"
	"errors"
	"fmt"
	"goauthz/security"

	"github.com/jmoiron/sqlx"
)

type SessionData struct {
	UserID int `db:"user_id"`
}

// Session ...
type Session struct {
	Data SessionData
	db   *sqlx.DB
	key  string
}

func (s *Session) Update() {
	fmt.Printf("Updating session %d\n", s.Data.UserID)
	if _, err := s.db.Exec("UPDATE sessions SET user_id = ? WHERE key = ?", s.Data.UserID, s.key); err != nil {
		panic(err)
	}
}

func createSession(db *sqlx.DB, key string, data SessionData) (*Session, error) {
	if _, err := db.Exec("INSERT INTO sessions (key, user_id) VALUES (?, ?)", key, data.UserID); err != nil {
		return nil, err
	}

	return &Session{
		db:   db,
		key:  key,
		Data: data,
	}, nil
}

func getSession(db *sqlx.DB, key string) (*Session, error) {
	var data SessionData
	if err := db.Get(&data, "SELECT user_id FROM sessions WHERE key = ?", key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: no session found", security.ErrUnauthenticated)
		}
		return nil, fmt.Errorf("%w: %v", security.ErrUnauthenticated, err)
	}

	return &Session{
		db:   db,
		key:  key,
		Data: data,
	}, nil
}
