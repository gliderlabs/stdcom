package console

import (
	"github.com/CloudyKit/jet"
	"github.com/gliderlabs/com/objects"
	"github.com/gliderlabs/stdcom/log"
	"github.com/gliderlabs/stdcom/web/auth"
)

func Register(registry *objects.Registry) {
	registry.Register(objects.New(&Component{}, ""))
}

type Component struct {
	Log           log.DebugLogger `com:"singleton"`
	Auth          auth.Requestor  `com:"singleton"`
	MenuProviders []MenuProvider  `com:"extpoint"`
	PageHandlers  []PageHandler   `com:"extpoint"`

	Views *jet.Set

	Config
}
