package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/vague2k/blkhell/internal/server/database"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userKey contextKey = "user"

type Service struct {
	db *database.Queries
}

func New(db *database.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) CreateNewUser(ctx context.Context, username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.CreateUser(ctx, database.CreateUserParams{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	})

	return err
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (*database.User, error) {
	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid username")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func (s *Service) CreateSession(ctx context.Context, userID int64) (string, time.Time, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", time.Time{}, err
	}
	sessionID := hex.EncodeToString(b)

	// 1 week before the session id expires
	expires := time.Now().Add(7 * (24 * time.Hour))

	_, err = s.db.CreateSession(ctx, database.CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expires,
	})

	return sessionID, expires, err
}

func (s *Service) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := s.db.GetSession(r.Context(), cookie.Value)
		if err != nil || time.Now().After(session.ExpiresAt) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.db.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
