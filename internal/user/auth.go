package user

import (
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

		if auth, ok := session.Values[userKey]; !ok || auth == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
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
