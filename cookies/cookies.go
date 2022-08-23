package cookies

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

const flashCookieName = "tuukoti_flash"

type Cookies struct {
	secureCookies *securecookie.SecureCookie
}

func New(sc *securecookie.SecureCookie) *Cookies {
	return &Cookies{secureCookies: sc}
}

func (c *Cookies) Create(name string, val string) (*http.Cookie, error) {
	encodedValue, err := c.secureCookies.Encode(name, val)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		Expires:  time.Now().Add(3 * (24 * time.Hour)),
		HttpOnly: true,
		Path:     "/",
	}

	return cookie, nil
}

// ReadCookie will decode the secure cookie value into `out`.
func (c *Cookies) Read(req *http.Request, name string, out interface{}) error {
	cookie, err := req.Cookie(name)
	if err != nil {
		return err
	}

	return c.secureCookies.Decode(name, cookie.Value, out)
}

// SetFlashCookie will set a flash cookie with the data you give it a map[string]string{}.
func SetFlash(w http.ResponseWriter, val map[string]string) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	co := &http.Cookie{
		Name:     flashCookieName,
		Value:    encode(b),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, co)

	return nil
}

// GetFlash will return any flash data as a map[string]string{}.
// It also sets the flash cookies value back to nothing so the flash only exists once.
func GetFlash(w http.ResponseWriter, req *http.Request) (map[string]string, error) {
	co, err := req.Cookie(flashCookieName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	value, err := decode(co.Value)
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	if err := json.Unmarshal(value, &out); err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     flashCookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Path:     "/",
		HttpOnly: true,
	})

	return out, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
