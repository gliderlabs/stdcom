package auth

import (
  "net/http"

  "github.com/markbates/goth"
)

type User = goth.User

type Requestor interface {
  CurrentUser(r *http.Request) *User
  LoginURL(r *http.Request, destPath string, providerName string) string
  LogoutURL(r *http.Request, destPath string) string
}
