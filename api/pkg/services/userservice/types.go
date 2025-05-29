package userservice

import (
	"fmt"
	"time"
)

type UserKind string

const (
	UserKindLocal   UserKind = "local"
	UserKindOIDC    UserKind = "oidc"
	UserKindService UserKind = "service"
)

var ErrUserNotFound = fmt.Errorf("unknown user")

type User struct {
	Id          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Admin       bool      `json:"admin"`
	Kind        UserKind  `json:"kind"`
	Active      bool      `json:"active"`
	Created     time.Time `json:"created"`
	LastLogin   time.Time `json:"last_login"`
}

type Session struct {
	Id      string    `json:"id"`
	Expires time.Time `json:"expires"`
}
