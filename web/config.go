package web

import "github.com/gliderlabs/com/config"

// Config is the configuration for Component
type Config struct {
	// Address and port to listen on
	ListenAddr string

	// Directory to serve static files from
	StaticDir string

	// URL path to serve static files at
	StaticPath string

	// Random string to use for session cookies
	CookieSecret string

	// Address and port to listen for TLS on
	TLSAddr string

	// Path to TLS cert file
	TLSCertPath string

	// Path to TLS key file
	TLSKeyPath string
}

// InitializeConfig sets up defaults and unmarshals into Config embedded in
// Component
func (c *Component) InitializeConfig(cfg config.Settings) error {
	cfg.SetDefault("ListenAddr", "0.0.0.0:8080")
	cfg.SetDefault("StaticDir", "ui/static/")
	cfg.SetDefault("StaticPath", "/static")
	cfg.SetDefault("TLSAddr", "0.0.0.0:4443")
	return cfg.Unmarshal(&(c.Config))
}
