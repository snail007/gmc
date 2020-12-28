package gcookie

import (
	gcore "github.com/snail007/gmc/core"
	"net/http"
	"time"
)

func New(w http.ResponseWriter, r *http.Request) (cookie *Cookies) {
	c := &Cookies{
		req: r,
		w:   w,
	}
	return c
}

type Cookies struct {
	req  *http.Request
	w    http.ResponseWriter
 }

func (c *Cookies) Get(name string) (value string, err error) {
	cookie, err := c.req.Cookie(name)
	if cookie == nil {
		return
	}
	value = cookie.Value
	return
}

// Set set the given cookie to the response and returns the current context to allow chaining.
// If options omit, it will use default options.
func (c *Cookies) Set(name, val string, options ...*gcore.CookieOptions) {
	opts := gcore.DefaultCookieOptions
	if len(options) > 0 {
		opts = options[0]
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    val,
		HttpOnly: opts.HTTPOnly,
		Secure:   opts.Secure,
		MaxAge:   opts.MaxAge,
		Domain:   opts.Domain,
		Path:     opts.Path,
	}
	if opts.MaxAge > 0 {
		d := time.Duration(opts.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d).UTC()
	} else if opts.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0).UTC()
	}
	http.SetCookie(c.w, cookie)
}

// Remove remove the given cookie
func (c *Cookies) Remove(name string, options ...*gcore.CookieOptions) {
	opts := *gcore.DefaultCookieOptions // should copy because we will change MaxAge
	if len(options) > 0 {
		opts = *options[0]
	}
	opts.MaxAge = -1
	c.Set(name, "", &opts)
}
