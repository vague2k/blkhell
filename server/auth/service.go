package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vague2k/blkhell/server/database"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userKey contextKey = "user"

var ErrDb = errors.New("Database error, try again. Contact admin if issue occurs")
var ErrNoSession = errors.New("Session does not exist")

type Service struct {
	db *database.Queries
}

func New(db *database.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) CreateNewUser(ctx context.Context, username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil && errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return errors.New("Password is too long") // err's if pass is longer than 72 bytes
	}

	_, err = s.db.CreateUser(ctx, database.CreateUserParams{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		return ErrDb
	}

	return nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (*database.User, error) {
	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid username")
		}
		return nil, ErrDb
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func (s *Service) CreateSession(ctx context.Context, userID string) (string, time.Time, error) {
	b := make([]byte, 32)
	rand.Read(b) // should never error, as it only returns error to fufill signature from io.Reader.Read()
	sessionToken := hex.EncodeToString(b)

	// 1 week before the session id expires
	expires := time.Now().Add(7 * (24 * time.Hour))

	_, err := s.db.CreateSession(ctx, database.CreateSessionParams{
		ID:        uuid.NewString(),
		Token:     sessionToken,
		UserID:    userID,
		ExpiresAt: expires,
	})
	if err != nil {
		return "", time.Time{}, errors.New("Database error, contact admin if issue occurs.")
	}

	return sessionToken, expires, err
}

func (s *Service) DestroySession(w http.ResponseWriter, r *http.Request) error {
	// can ONLY return error if no cookie is found
	// if no cookie is found on logout, user should already logged out ???
	// hmm..
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	err = s.db.DeleteSession(r.Context(), cookie.Value)
	if err != nil {
		return ErrDb
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1, // MaxAge<0 means delete cookie now
	})

	return nil
}

func (s *Service) UserFromContext(ctx context.Context) (*database.User, bool) {
	u, ok := ctx.Value(userKey).(*database.User)
	return u, ok
}

func (s *Service) GetUserFromRequest(r *http.Request) (*database.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, ErrNoSession
	}

	session, err := s.db.GetSessionByToken(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSession
		}
		return nil, ErrDb
	}

	user, err := s.db.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Couldn't find user tied to session")
		}
		return nil, ErrDb
	}

	return &user, nil
}

func (s *Service) RedirectIfAuth(next http.Handler) http.Handler {
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

func (s *Service) RequireAuth(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), userKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
