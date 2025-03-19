// Package zauth provides auth integration with the Zeus authentication system
package zauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const endpoint = "https://zauth.zeus.gent"

// Provider is the implementation of `goth.Provider` for accessing Zauth
type Provider struct {
	config     *oauth2.Config
	httpClient *http.Client

	clientKey    string
	secret       string
	callbackURL  string
	userURL      string
	providerName string
}

// Interface compliance
var _ goth.Provider = (*Provider)(nil)

// User contains the user data received from Zauth
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	FullName string `json:"full_name"`
}

// New creates a new Zauth provider
func New(clientKey, secret, callbackURL string) *Provider {
	p := &Provider{
		clientKey:    clientKey,
		secret:       secret,
		callbackURL:  callbackURL,
		userURL:      fmt.Sprintf("%s/current_user", endpoint),
		providerName: "zauth",
	}
	c := &oauth2.Config{
		ClientID:     p.clientKey,
		ClientSecret: p.secret,
		RedirectURL:  p.callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/oauth/authorize", endpoint),
			TokenURL:  fmt.Sprintf("%s/oauth/token", endpoint),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
	p.config = c

	return p
}

// Interface methods

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// BeginAuth asks Zauth for an authentication endpoint.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	s := &Session{
		AuthURL: url,
	}

	return s, nil
}

// UnmarshalSession a string into a session
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}

	err := json.NewDecoder(strings.NewReader(data)).Decode(s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data into session %s | %v", data, err)
	}

	return s, nil
}

// FetchUser will go to Zauth and access basic information about the user.
func (p *Provider) FetchUser(gothSession goth.Session) (goth.User, error) {
	s := gothSession.(*Session)
	if s.AccessToken == "" {
		return goth.User{}, fmt.Errorf("unable to fetch user information without an access token %+v", *s)
	}

	req, err := http.NewRequestWithContext(goth.ContextForClient(p.client()), "GET", p.userURL, nil)
	if err != nil {
		return goth.User{}, fmt.Errorf("unable to create a new http request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	response, err := p.client().Do(req)
	if err != nil {
		return goth.User{}, fmt.Errorf("received error from oauth2 user fetch call %+v | %v", *p, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return goth.User{}, fmt.Errorf("received wrong http status code %d | %+v | %v", response.StatusCode, *p, err)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return goth.User{}, fmt.Errorf("unable to read zauth response body %v", err)
	}

	var u User
	if err := json.Unmarshal(responseBytes, &u); err != nil {
		return goth.User{}, fmt.Errorf("unable to unmarshal zauth response into zauthUser %s | %v", string(responseBytes), err)
	}

	user := goth.User{
		AccessToken:  s.AccessToken,
		Provider:     p.Name(),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,

		RawData: map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"fullName": u.FullName,
			"admin":    u.Admin,
		},
	}

	return user, nil
}

// Debug is a no-op for Zauth
func (p *Provider) Debug(_ bool) {}

// RefreshToken get new access token based on the refresh token
// Not implemented for Zauth
func (p *Provider) RefreshToken(_ string) (*oauth2.Token, error) {
	return nil, errors.New("zauth doesn't support refresh tokens")
}

// RefreshTokenAvailable returns if refresh tokens are supported
// This is not the case for  Zauth
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}

// Private methods

func (p *Provider) client() *http.Client {
	return goth.HTTPClientWithFallBack(p.httpClient)
}
