// Package init is used to register the sessions component with the default registry
// via side effect import.
package init

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/web/sessions"
)

func init() {
	sessions.Register(com.DefaultRegistry)
}
