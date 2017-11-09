package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/markbates/goth"
)

const (
	SessionKeyFmt    = "_auth_%s"
	ReturnQueryParam = "return"
)

var (
	webRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	UserSessionKey     = fmt.Sprintf(SessionKeyFmt, "user")
	ProviderSessionKey = fmt.Sprintf(SessionKeyFmt, "provider")
	ReturnSessionKey   = fmt.Sprintf(SessionKeyFmt, "return")
)

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, c.BasePath)
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc(c.absPath("/login/"), c.handleLogin)
	mux.HandleFunc(c.absPath("/callback/"), c.handleCallback)
	mux.HandleFunc(c.absPath("/logout"), c.handleLogout)
	mux.ServeHTTP(w, r)
}

func (c *Component) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := c.Session.Unset(w, r, UserSessionKey)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	providerName := c.Session.ValueString(r, ProviderSessionKey)
	err = c.Session.Unset(w, r, fmt.Sprintf(SessionKeyFmt, providerName))
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	redirect := r.URL.Query().Get(ReturnQueryParam)
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func (c *Component) handleLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get(ReturnQueryParam)
	if redirect == "" {
		redirect = r.Referer()
	}
	err := c.Session.Set(w, r, ReturnSessionKey, redirect)
	if err != nil {
		c.Log.Debug(err)
	}

	providerName := path.Base(r.URL.Path)
	url, err := c.setupProviderURL(w, r, providerName)
	if err != nil {
		c.httpError(w, err, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *Component) handleCallback(w http.ResponseWriter, r *http.Request) {
	providerName := path.Base(r.URL.Path)
	sessionKey := fmt.Sprintf(SessionKeyFmt, providerName)

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		c.httpError(w, err, http.StatusNotFound)
		return
	}

	value := c.Session.ValueString(r, sessionKey)
	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	err = validateState(r, sess)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	var u User
	u, err = provider.FetchUser(sess)
	if err != nil {
		// FetchUser failed, so try again:
		// get new token and retry fetch
		_, err = sess.Authorize(provider, r.URL.Query())
		if err != nil {
			c.httpError(w, err, http.StatusInternalServerError)
			return
		}

		err = c.Session.Set(w, r, sessionKey, sess.Marshal())
		if err != nil {
			c.httpError(w, err, http.StatusInternalServerError)
			return
		}

		u, err = provider.FetchUser(sess)
		if err != nil {
			c.httpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	err = c.Session.Set(w, r, UserSessionKey, u)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	err = c.Session.Set(w, r, ProviderSessionKey, providerName)
	if err != nil {
		c.httpError(w, err, http.StatusInternalServerError)
		return
	}

	// TODO: tell listeners

	redirect := c.Session.ValueString(r, ReturnSessionKey)
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func (c *Component) setupProviderURL(w http.ResponseWriter, r *http.Request, providerName string) (string, error) {
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(setupState(r))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = c.Session.Set(w, r, fmt.Sprintf(SessionKeyFmt, providerName), sess.Marshal())
	if err != nil {
		return "", err
	}

	return url, err
}

func (c *Component) absPath(parts ...string) string {
	return strings.Join(append([]string{c.BasePath}, parts...), "")
}

func (c *Component) httpError(w http.ResponseWriter, err error, status int) {
	http.Error(w, http.StatusText(status), status)
	c.Log.Debug(err)
}

func setupState(r *http.Request) string {
	state := r.URL.Query().Get("state")
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(webRand.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func validateState(req *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != req.URL.Query().Get("state")) {
		return errors.New("state token mismatch")
	}
	return nil
}
