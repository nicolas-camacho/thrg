package user

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nicolas-camacho/thrg/internal/contextutil"
)

const (
	SessionName = "admin_session"
	userKey     = "user_id"
)

var Store *sessions.CookieStore

func login(w http.ResponseWriter, r *http.Request, userID uuid.UUID, sessionName string) error {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("error retrieving session: %w", err)
	}

	session.Values[userKey] = userID

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}
	return nil
}

func LoginAdmin(w http.ResponseWriter, r *http.Request, userID uuid.UUID, adminSessionName string) error {
	return login(w, r, userID, adminSessionName)
}

func LoginPlayer(w http.ResponseWriter, r *http.Request, userID uuid.UUID, playerSessionName string) error {
	return login(w, r, userID, playerSessionName)
}

func LogoutUser(w http.ResponseWriter, r *http.Request, sessionName string) error {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return nil
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}
	return nil
}

func AdminAuthMiddleware(adminSessionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := Store.Get(r, adminSessionName)
			if err != nil {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			userID, ok := session.Values[userKey].(uuid.UUID)
			if !ok || userID == uuid.Nil {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			ctx := contextutil.SetUserIDInContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PlayerAuthMiddleware(playerSessionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := Store.Get(r, playerSessionName)
			if err != nil || session.IsNew {
				http.Redirect(w, r, "/player/login", http.StatusSeeOther)
				return
			}

			userID, ok := session.Values[userKey].(uuid.UUID)
			if !ok || userID == uuid.Nil {
				http.Redirect(w, r, "/player/login", http.StatusSeeOther)
				return
			}

			// Agrega el ID del jugador al contexto
			ctx := contextutil.SetUserIDInContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
