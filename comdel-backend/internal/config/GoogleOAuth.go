package config

import (
	"context"

	"golang.org/x/oauth2"
)

type GoogleOAuth struct {
	Conf *oauth2.Config
} 

func NewGoogleOAuth(conf *oauth2.Config) OAuthProvider {
	return &GoogleOAuth{Conf:  conf}
}

func (goa *GoogleOAuth) TokenSource(ctx context.Context ,token *oauth2.Token) oauth2.TokenSource {
	return goa.Conf.TokenSource(ctx, token)
}