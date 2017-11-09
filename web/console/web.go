package console

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/CloudyKit/jet"
)

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, c.BasePath)
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.TrimPrefix(r.URL.Path, c.BasePath) {
	case "/login":
		if c.Auth.CurrentUser(r) != nil {
			http.Redirect(w, r, c.absPath("/"), http.StatusTemporaryRedirect)
			return
		}
		c.RenderTemplate(w, r, "login.jet", nil, nil)
	default:
		if c.Auth.CurrentUser(r) == nil {
			http.Redirect(w, r, c.absPath("/login"), http.StatusTemporaryRedirect)
			return
		}
		ctx := context.WithValue(r.Context(), "auth", c.Auth)
		req := r.WithContext(ctx)
		for _, handler := range c.PageHandlers {
			if handler.MatchHTTP(req) {
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, r)
				if loc := rr.HeaderMap.Get("Location"); loc != "" {
					http.Redirect(w, r, loc, http.StatusTemporaryRedirect)
					return
				}
				c.RenderTemplate(w, r, "page.jet", nil, rr.Body.String())
				return
			}
		}
		c.RenderTemplate(w, r, "index.jet", nil, nil)
	}
}

func (c *Component) absPath(parts ...string) string {
	return strings.Replace(strings.Join(append([]string{c.BasePath}, parts...), ""), "//", "/", -1)
}

func (c *Component) httpError(w http.ResponseWriter, err error, status int) {
	http.Error(w, http.StatusText(status), status)
	c.Log.Debug(err)
}

func (c *Component) RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, vars map[string]interface{}, data interface{}) {
	t, err := c.Views.GetTemplate(tmpl)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}
	v := make(jet.VarMap)
	v.Set("req", r)
	v.Set("auth", c.Auth)
	v.Set("user", c.Auth.CurrentUser(r))
	var items []MenuItem
	for _, provider := range c.MenuProviders {
		for _, item := range provider.ConsoleMenuItems() {
			items = append(items, item)
		}
	}
	v.Set("menu", items)
	if vars != nil {
		for k := range vars {
			v.Set(k, vars[k])
		}
	}
	if err = t.Execute(w, v, data); err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
	}
}
