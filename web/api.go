package web

import (
	"net/http"
)

// Handler extension API for matching and handling HTTP requests
type Handler interface {
	// MatchHTTP returns true when a Request should be handled by this object
	MatchHTTP(r *http.Request) bool

	// ServeHTTP handles the Request if MatchHTTP returns true
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Middleware handler API
type Middleware func(http.Handler) http.Handler

// MiddlewareInjector defines ordered HTTP middleware to use
type MiddlewareInjector interface {
	// Middleware returns a list of middleware in the order of outer first
	Middleware() []Middleware
}
