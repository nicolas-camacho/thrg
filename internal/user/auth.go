package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "admin_session"
	userKey     = "user_id"
)

var Store *sessions.CookieStore

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, SessionName)
		if err != nil {
			log.Printf("Error retrieving session: %v", err)
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}

		authID, ok := session.Values[userKey]
		if !ok || authID == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}

		userID, ok := authID.(uint)
		if !ok || userID == 0 {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return fmt.Errorf("error retrieving session: %w", err)
	}

	session.Values[userKey] = userID
	session.Values["role"] = RoleAdmin

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return nil
}
