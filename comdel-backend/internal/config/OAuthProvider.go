package config

import (
	"context"

	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	TokenSource(context.Context ,*oauth2.Token) oauth2.TokenSource
}