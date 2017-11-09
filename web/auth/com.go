package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gliderlabs/com/objects"
	"github.com/gliderlabs/stdcom/log"
	"github.com/gliderlabs/stdcom/web/sessions"
)

func Register(reg *objects.Registry) {
	reg.Register(&objects.Object{Value: &Component{}})
}

type Component struct {
	Log     log.DebugLogger  `com:"singleton"`
	Session sessions.Session `com:"singleton"`

	Config
}

func (c *Component) CurrentUser(r *http.Request) *User {
	v := c.Session.Value(r, fmt.Sprintf(SessionKeyFmt, "user"))
	if v == nil {
		return nil
	}
	u, ok := v.(User)
	if !ok {
		return nil
	}
	return &u
}

func (c *Component) LoginURL(r *http.Request, providerName string, dest string) string {
	q := url.Values{}
	if dest != "" {
		q.Set(ReturnQueryParam, dest)
	} else {
		q.Set(ReturnQueryParam, r.URL.String())
	}
	return fmt.Sprintf("%s?%s", c.absPath("/login/", providerName), q.Encode())
}

func (c *Component) LogoutURL(r *http.Request, dest string) string {
	q := url.Values{}
	if dest != "" {
		q.Set(ReturnQueryParam, dest)
	} else {
		q.Set(ReturnQueryParam, r.URL.String())
	}
	return fmt.Sprintf("%s?%s", c.absPath("/logout"), q.Encode())
}
