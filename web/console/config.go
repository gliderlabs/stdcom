package console

import (
	"github.com/gliderlabs/com/config"
)

type Config struct {
	// base path to use for web endpoints
	BasePath string

	// name of console used in titles
	Name string
}

func (c *Component) InitializeConfig(cfg config.Settings) error {
	cfg.SetDefault("BasePath", "/console")
	return cfg.Unmarshal(&(c.Config))
}
