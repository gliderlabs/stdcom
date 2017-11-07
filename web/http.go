package web

import (
	"context"
	"net/http"
)

// ServeHTTP wraps the Handler protocol and injects middleware
func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler

	// set up a handler that implements the core match+handle protocol
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range c.Handlers {
			if handler.MatchHTTP(r) {
				handler.ServeHTTP(w, r)
				return
			}
		}

		// attempt handlers again allowing catch
		// of NotFound via URL Fragment
		r.URL.Fragment = "NotFound"
		for _, handler := range c.Handlers {
			if handler.MatchHTTP(r) {
				handler.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)
	})

	// unwrap the middleware handlers
	// if middleware injector is present
	if c.Middleware != nil {
		middleware := c.Middleware.Middleware()
		for i := range middleware {
			h = middleware[len(middleware)-1-i](h)
		}
	}

	// serve with the component instance in the request context
	ctx := context.WithValue(r.Context(), ComponentContextKey, c)
	h.ServeHTTP(w, r.WithContext(ctx))
}
