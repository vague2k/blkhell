package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/server/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	config *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{config: cfg}
}

func (s *AuthService) CreateNewUser(ctx context.Context, username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil && errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return errors.New("Password is too long") // err's if pass is longer than 72 bytes
	}

	_, err = s.config.Database.CreateUser(ctx, database.CreateUserParams{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		return serverErrors.ErrDb
	}

	return nil
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (*database.User, error) {
	user, err := s.config.Database.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid username")
		}
		return nil, serverErrors.ErrDb
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func (s *AuthService) CreateSession(ctx context.Context, userID string) (string, time.Time, error) {
	b := make([]byte, 32)
	rand.Read(b) // should never error, as it only returns error to fufill signature from io.Reader.Read()
	sessionToken := hex.EncodeToString(b)

	// 1 week before the session id expires
	expires := time.Now().Add(7 * (24 * time.Hour))

	_, err := s.config.Database.CreateSession(ctx, database.CreateSessionParams{
		ID:        uuid.NewString(),
		Token:     sessionToken,
		UserID:    userID,
		ExpiresAt: expires,
	})
	if err != nil {
		return "", time.Time{}, serverErrors.ErrDb
	}

	return sessionToken, expires, nil
}

func (s *AuthService) DestroySession(w http.ResponseWriter, r *http.Request) error {
	// can ONLY return error if no cookie is found
	// if no cookie is found on logout, user should already logged out ???
	// hmm..
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	err = s.config.Database.DeleteSession(r.Context(), cookie.Value)
	if err != nil {
		return serverErrors.ErrDb
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

func (s *AuthService) UserFromContext(ctx context.Context) (*database.User, bool) {
	u, ok := ctx.Value(middleware.AuthUserKey).(*database.User)
	return u, ok
}
