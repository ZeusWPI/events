package zauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/markbates/goth"
)

// Session is the implementation of `goth.Session` for storing data during the auth process with Zauth
type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// Interface compliance
var _ goth.Session = (*Session)(nil)

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Zauth provider.
func (s *Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Marshal the session into a string
func (s *Session) Marshal() string {
	b, _ := json.Marshal(*s)
	return string(b)
}

// Authorize the session with Zauth and return the access token to be stored for future use
func (s *Session) Authorize(gothProvider goth.Provider, params goth.Params) (string, error) {
	p := gothProvider.(*Provider)
	token, err := p.config.Exchange(goth.ContextForClient(p.client()), params.Get("code"))
	if err != nil {
		return "", fmt.Errorf("unable to exchange codes %v", err)
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.ExpiresAt = token.Expiry

	return token.AccessToken, nil
}
