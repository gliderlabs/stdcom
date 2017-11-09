package console

import (
	"path"
	"runtime"
	"time"

	"github.com/CloudyKit/jet"
	humanize "github.com/dustin/go-humanize"
)

func (c *Component) InitializeDaemon() error {
	_, filename, _, _ := runtime.Caller(0)
	c.Views = jet.NewHTMLSet(path.Join(path.Dir(filename), "html"))
	c.Views.SetDevelopmentMode(true)
	c.Views.AddGlobal("config", c.Config)
	c.Views.AddGlobal("time", func(t interface{}) string {
		return humanize.Time(t.(time.Time))
	})
	return nil
}
