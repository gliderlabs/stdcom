// Package init is used to register the auth component with the default registry
// via side effect import.
package init

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/web/auth"
)

func init() {
	auth.Register(com.DefaultRegistry)
}
