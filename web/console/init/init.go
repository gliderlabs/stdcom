// Package init is used to register the console component with the default registry
// via side effect import.
package init

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/web/console"
)

func init() {
	console.Register(com.DefaultRegistry)
}
