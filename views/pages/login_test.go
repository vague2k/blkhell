package pages_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vague2k/blkhell/testutil"
	"github.com/vague2k/blkhell/views/pages"
)

func TestLogin(t *testing.T) {

	t.Run("UI renders page title", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		assert.Equal(t, "blkhell", doc.Find("h1").First().Text())
	})

	t.Run("UI renders subtitle", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		assert.Equal(t, "For Blackheaven members only...", doc.Find("p").First().Text())
	})

	t.Run("UI renders login form with correct attributes", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		form := doc.Find("form")
		assert.Equal(t, 1, form.Length())
		assert.Equal(t, "/login", form.AttrOr("hx-post", ""))
		assert.Equal(t, "#global-toast", form.AttrOr("hx-target", ""))
	})

	t.Run("UI renders username input", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		input := doc.Find("input[name=\"username\"]")
		assert.Equal(t, 1, input.Length())
		assert.Equal(t, "text", input.AttrOr("type", ""))
	})

	t.Run("UI renders password input", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		input := doc.Find("input[name=\"password\"]")
		assert.Equal(t, 1, input.Length())
		assert.Equal(t, "password", input.AttrOr("type", ""))
	})

	t.Run("UI renders login button", func(t *testing.T) {
		test := testutil.NewTest(t)
		doc, err := test.RenderComponent(pages.Login())
		require.NoError(t, err)
		button := doc.Find("button[type=\"submit\"]")
		assert.Equal(t, 1, button.Length())
		assert.Equal(t, "Login", button.Text())
	})

	t.Run("Handler redirects to dashboard with valid credentials", func(t *testing.T) {
		test := testutil.NewTest(t)

		username := test.RandomUsername()
		password := test.RandomPassword()
		err := test.AuthService.CreateNewUser(test.Context(), username, password, "admin")
		require.NoError(t, err)

		form := url.Values{}
		form.Add("username", username)
		form.Add("password", password)

		req := test.NewFormRequest("/login", form)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := test.NewRecorder()

		test.Handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "/dashboard", w.Header().Get("HX-Redirect"))

		cookie := w.Result().Cookies()
		assert.Len(t, cookie, 1)
		assert.Equal(t, "session_token", cookie[0].Name)
	})

	t.Run("Handler returns error with invalid username", func(t *testing.T) {
		test := testutil.NewTest(t)

		form := url.Values{}
		form.Add("username", test.RandomUsername())
		form.Add("password", test.RandomPassword())

		req := test.NewFormRequest("/login", form)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := test.NewRecorder()

		test.Handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "invalid username")
	})

	t.Run("Handler returns error with invalid password", func(t *testing.T) {
		test := testutil.NewTest(t)

		username := test.RandomUsername()
		err := test.AuthService.CreateNewUser(test.Context(), username, test.RandomPassword(), "admin")
		require.NoError(t, err)

		form := url.Values{}
		form.Add("username", username)
		form.Add("password", test.RandomPassword())

		req := test.NewFormRequest("/login", form)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := test.NewRecorder()

		test.Handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "invalid password")
	})
}
