package console

import "github.com/gliderlabs/stdcom/web"

type MenuProvider interface {
	ConsoleMenuItems() []MenuItem
}

type MenuItem struct {
	Title string
	Link  string
}

type PageHandler interface {
	MenuProvider
	web.Handler
}
