package auth

import (
	"github.com/gliderlabs/com/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

type Config struct {
	// base path to use for web endpoints
	BasePath string

	// TEMP
	ClientKey    string
	ClientSecret string
}

func (c *Component) InitializeConfig(cfg config.Settings) (err error) {
	cfg.SetDefault("BasePath", "/_auth")
	err = cfg.Unmarshal(&(c.Config))
	// TODO: REMOVE
	goth.UseProviders(
		github.New(c.ClientKey, c.ClientSecret, "http://localhost:8080/_auth/callback/github"))

	return
}
