package cookies_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/require"
	"github.com/tuukoti/framework/cookies"
)

func TestCookieReadAndGet(t *testing.T) {
	secureCookie := securecookie.New([]byte("very-secret"), []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	c := cookies.New(secureCookie)

	w := httptest.NewRecorder()

	cookie, err := c.Create("testing", "test")
	require.NoError(t, err)

	http.SetCookie(w, cookie)

	req := &http.Request{
		Header: http.Header{},
	}

	req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	out := ""

	err = c.Read(req, "testing", &out)
	require.NoError(t, err)
	require.Equal(t, "test", out)
}

func TestSetGetFlash(t *testing.T) {
	w := httptest.NewRecorder()

	d := map[string]string{
		"test": "testing",
	}

	err := cookies.SetFlash(w, d)
	require.NoError(t, err)

	req := &http.Request{
		Header: http.Header{},
	}

	req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	data, err := cookies.GetFlash(w, req)
	require.NoError(t, err)
	require.Equal(t, d, data)
}
