package sessions

import (
  "net/http"

  "github.com/gliderlabs/com/objects"
  "github.com/gorilla/sessions"
)

func Register(reg *objects.Registry) {
	reg.Register(&objects.Object{Value: &Component{}})
}

type Store = sessions.Store

type Component struct {
  Initializer Initializer `com:"singleton"`

	store Store
}

type Initializer interface {
  InitializeSessions() (Store, error)
}

type Session interface {
  Value(r *http.Request, key string) interface{}
  ValueString(r *http.Request, key string) string
  Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error
  Unset(w http.ResponseWriter, r *http.Request, key string) error
}

func (c *Component) InitializeDaemon() error {
  return c.Initialize()
}

func (c *Component) Initialize() (err error) {
  if c.Initializer == nil {
    // TODO: make configurable
    // key := securecookie.GenerateRandomKey(64)
    // if key == nil {
    //   return errors.New("nil random key generation")
    // }
    cs := sessions.NewCookieStore([]byte("replace-me"))
    cs.Options.HttpOnly = true
    c.store = cs
    return nil
  }
  c.store, err = c.Initializer.InitializeSessions()
  return
}

func (c *Component) Get(r *http.Request, name string) (*sessions.Session, error) {
  return c.store.Get(r, name)
}

func (c *Component) New(r *http.Request, name string) (*sessions.Session, error) {
  return c.store.New(r, name)
}

func (c *Component) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
  return c.store.Save(r, w, s)
}

func (c *Component) Value(r *http.Request, key string) interface{} {
  session, err := c.store.Get(r, key)
	if err != nil {
		return nil
	}
	val, exists := session.Values[key]
	if !exists {
		return nil
	}
	return val
}

func (c *Component) ValueString(r *http.Request, key string) string {
  s, _ := c.Value(r, key).(string)
  return s
}

func (c *Component) Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
  session, err := c.store.Get(r, key)
	if err != nil {
		return err
	}
	session.Values[key] = value
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func (c *Component) Unset(w http.ResponseWriter, r *http.Request, key string) error {
  session, err := c.store.Get(r, key)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}
