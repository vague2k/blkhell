package pages_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vague2k/blkhell/server/handlers"
	"github.com/vague2k/blkhell/server/services"
	"github.com/vague2k/blkhell/testutil"
	"github.com/vague2k/blkhell/views/pages"
)

func TestLogin(t *testing.T) {
	t.Run("UI renders page title", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		assert.Equal(t, "blkhell", doc.Find("h1").First().Text())
	})

	t.Run("UI renders subtitle", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		assert.Equal(t, "For Blackheaven members only...", doc.Find("p").First().Text())
	})

	t.Run("UI renders login form with correct attributes", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		form := doc.Find("form")
		assert.Equal(t, 1, form.Length())
		assert.Equal(t, "/login", form.AttrOr("hx-post", ""))
		assert.Equal(t, "#global-toast", form.AttrOr("hx-target", ""))
	})

	t.Run("UI renders username input", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		input := doc.Find("input[name=\"username\"]")
		assert.Equal(t, 1, input.Length())
		assert.Equal(t, "text", input.AttrOr("type", ""))
	})

	t.Run("UI renders password input", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		input := doc.Find("input[name=\"password\"]")
		assert.Equal(t, 1, input.Length())
		assert.Equal(t, "password", input.AttrOr("type", ""))
	})

	t.Run("UI renders login button", func(t *testing.T) {
		doc, err := testutil.RenderComponent(pages.Login())
		require.NoError(t, err)
		button := doc.Find("button[type=\"submit\"]")
		assert.Equal(t, 1, button.Length())
		assert.Equal(t, "Login", button.Text())
	})

	t.Run("Handler redirects to dashboard with valid credentials", func(t *testing.T) {
		cfg, cleanup := testutil.NewTestConfig(t)
		t.Cleanup(cleanup)

		auth := services.NewAuthService(cfg)
		handler := handlers.NewHandler(cfg)

		err := auth.CreateNewUser(context.Background(), "testuser", "password123", "admin")
		require.NoError(t, err)

		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "/dashboard", w.Header().Get("HX-Redirect"))

		cookie := w.Result().Cookies()
		assert.Len(t, cookie, 1)
		assert.Equal(t, "session_token", cookie[0].Name)
	})

	t.Run("Handler returns error with invalid username", func(t *testing.T) {
		cfg, cleanup := testutil.NewTestConfig(t)
		t.Cleanup(cleanup)

		handler := handlers.NewHandler(cfg)

		form := url.Values{}
		form.Add("username", "nonexistent")
		form.Add("password", "password123")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "invalid username")
	})

	t.Run("Handler returns error with invalid password", func(t *testing.T) {
		cfg, cleanup := testutil.NewTestConfig(t)
		t.Cleanup(cleanup)

		auth := services.NewAuthService(cfg)
		handler := handlers.NewHandler(cfg)

		err := auth.CreateNewUser(context.Background(), "testuser2", "password123", "admin")
		require.NoError(t, err)

		form := url.Values{}
		form.Add("username", "testuser2")
		form.Add("password", "wrongpassword")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "invalid password")
	})
}
