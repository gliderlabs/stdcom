package web

import (
	"net/http"
	"strings"
)

// ForceSecureMiddleware will redirect requests to their HTTPS equivalent if
// HTTPS has been configured for the component.
func ForceSecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := FromRequest(r)

		// redirect to https if available
		if r.TLS == nil && c.https != nil {
			u := r.URL
			u.Host = r.Host
			u.Scheme = "https"
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StaticFileMiddleware serves static files based on component configuration of
// StaticPath and StaticDir.
func StaticFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := FromRequest(r)

		// serve static
		staticPrefix := c.StaticPath + "/"
		if strings.HasPrefix(r.URL.Path, staticPrefix) {
			http.StripPrefix(staticPrefix,
				http.FileServer(http.Dir(c.StaticDir))).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
