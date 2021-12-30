package app

import (
	"context"
	"goauthz/security"
	"goauthz/usersession"

	"github.com/jmoiron/sqlx"
)

// App ...
type App struct {
	DB    *sqlx.DB
	auth  *security.AuthenticationRegistry
	authz *security.AuthorizationRegistry

	sessions *usersession.SessionManager
}

func New(db *sqlx.DB) *App {
	app := &App{
		DB: db,
	}

	app.registerSecurityMethods()
	return app
}

func (app *App) registerSecurityMethods() {
	// Create the user session manager
	// We pass a principal-creator function that will create a principal from the session data
	app.sessions = usersession.New(app.DB, func(c context.Context, session *usersession.Session) (security.Principal, error) {
		return &User{
			UserID: int(session.Data.UserID),
			Rights: make([]string, 0),
		}, nil
	})

	// Register the session authenticator to the usersession credential type
	// Whenever a credential is provided with the type `usersession.Type`, this session authenticator will be used
	app.auth = security.NewAuthenticationRegistery().
		Add(usersession.Type, app.sessions)

	// Register the user principal authorizer to the user principal type
	// Notice in the usersession constructor we provide a principal-creator function that will create a principal from the session data
	// this will create the user principal which will be used by this authorizer
	app.authz = security.NewAuthorizerRegistry().
		Add(UserPrincipalType, NewUserAuthorizer())
}
