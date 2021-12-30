package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"goauthz/security"
	"goauthz/usersession"

	"github.com/go-chi/chi"
)

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	// Auth middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Extract user id from authorization header
			hdr := r.Header.Get("authorization")
			if hdr == "" {
				next.ServeHTTP(rw, r)
				return
			}

			parts := strings.Split(hdr, " ")

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				sendErr(rw, security.ErrUnauthenticated)
				return
			}

			// Get credential
			credential, err := security.UnmarshalCredentialString(parts[1])
			if err != nil {
				fmt.Printf("Error while unmarshalling credential string: %v\n", err)
				sendErr(rw, security.ErrUnauthenticated)
				return
			}

			authCTX, err := a.auth.Authenticate(r.Context(), credential)
			if err != nil {
				fmt.Printf("Error while authenticating credential: %v\n", err)
				sendErr(rw, security.ErrUnauthenticated)
				return
			}

			// Continue
			next.ServeHTTP(rw, r.WithContext(authCTX))
		})
	})

	r.Post("/login", a.httpPostLogin())
	r.Get("/users/{uid}", a.httpGetUserData())

	return r
}

func (app *App) httpPostLogin() http.HandlerFunc {
	// loginDTO ...
	type loginDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		var dto loginDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			sendErr(rw, err)
			return
		}

		// Fetch the user with the given username and password
		// obviously never store passwords in plain text!
		var userID int
		if err := app.DB.Get(&userID, "SELECT id FROM users WHERE username = ? AND password = ?", dto.Username, dto.Password); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				sendErr(rw, security.ErrUnauthenticated)
				return
			}
			sendErr(rw, err)
			return
		}

		// Create session
		_, credentials, err := app.sessions.CreateSession(usersession.SessionData{UserID: userID})
		if err != nil {
			sendErr(rw, err)
			return
		}

		sendJSON(rw, map[string]interface{}{"message": "login successful", "data": credentials.String()})
	}
}

func (app *App) httpGetUserData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Get requested user ID
		uidStr := chi.URLParam(r, "uid")
		uid, err := strconv.Atoi(uidStr)
		if err != nil {
			sendErr(rw, err)
			return
		}

		if err := app.authz.MustHaveRights(r.Context(), security.UserIDToURN(uid), "user.read"); err != nil {
			if errors.Is(err, security.ErrForbidden) {
				sendErr(rw, errors.New("not found"))
				return
			}
			sendErr(rw, err)
			return
		}

		sendJSON(rw, map[string]interface{}{
			"user_id": uid,
			"access":  "granted!",
		})
	}
}

func sendErr(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func sendJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(rw).Encode(data); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
