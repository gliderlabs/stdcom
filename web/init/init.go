// Package init is used to register the web component with the default registry
// via side effect import.
package init

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/web"
)

func init() {
	web.Register(com.DefaultRegistry)
}
