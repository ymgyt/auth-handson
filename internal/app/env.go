package app

import (
	"firebase.google.com/go/auth"
)

// Env -
type Env struct {
	AppRoot string
	Auth    *auth.Client

	// Auth0
	Auth0 Auth0Env
}

type Auth0Env struct {
	Domain       string
	ClientSecret string
	ClientID     string

	BaseURL string
}
