package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/vague2k/blkhell/config"
)

type ctxKey string

const AuthUserKey ctxKey = "user"

type Middleware struct {
	config *config.Config
}

func New(cfg *config.Config) *Middleware {
	return &Middleware{config: cfg}
}

func (m *Middleware) RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := m.config.Database.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil || time.Now().After(session.ExpiresAt) {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := m.config.Database.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil || time.Now().After(session.ExpiresAt) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := m.config.Database.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
