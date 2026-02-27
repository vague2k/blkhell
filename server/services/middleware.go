package services

import (
	"context"
	"net/http"
	"time"

	"github.com/vague2k/blkhell/server/database"
)

type ctxKey string

const authUserKey ctxKey = "user"

type MiddlewareService struct {
	db *database.Queries
}

func NewMiddlewareService(db *database.Queries) *MiddlewareService {
	return &MiddlewareService{db: db}
}

func (s *MiddlewareService) RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := s.db.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil || time.Now().After(session.ExpiresAt) {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *MiddlewareService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := s.db.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil || time.Now().After(session.ExpiresAt) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.db.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *MiddlewareService) Bands(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bands, err := s.db.GetBands(r.Context())
		if err != nil {
			http.Error(w, "could not get bands", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), bandsCtxKey, bands)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
