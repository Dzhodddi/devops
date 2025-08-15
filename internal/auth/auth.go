package auth

import (
	"devops/internal/env"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"strings"
)

const (
	maxAge = 60 * 60 * 24 * 7 // 7 days
	isProd = false
)

var (
	clientID     = env.GetString("CLIENT_ID", "")
	clientSecret = env.GetString("CLIENT_SECRET", "")
	key          = env.GetString("SECRET_KEY", "")
)

func NewAuth() {
	goth.UseProviders(
		google.New(clientID, clientSecret, "http://localhost:3000/auth/google/callback"),
	)

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Domain = ""
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

}

func BeginAuth(w http.ResponseWriter, r *http.Request, provider string) {
	q := r.URL.Query()
	q.Set("provider", strings.ToLower(provider))
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func CompleteAuth(w http.ResponseWriter, r *http.Request, provider string) (goth.User, error) {
	q := r.URL.Query()
	q.Set("provider", strings.ToLower(provider))
	r.URL.RawQuery = q.Encode()
	return gothic.CompleteUserAuth(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	return gothic.Logout(w, r)
}

func GetUserFromSession(r *http.Request) (*goth.User, error) {
	// This will look up the provider name stored in the session
	provider, err := gothic.GetFromSession("provider", r)
	if err != nil {
		return nil, err // no provider in session → not logged in
	}

	// Get the provider instance
	p, err := goth.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	// Get the existing session
	sess, err := gothic.GetFromSession(provider, r)
	if err != nil {
		return nil, err
	}

	// Unmarshal into provider's session type
	s, err := p.UnmarshalSession(sess)
	if err != nil {
		return nil, err
	}

	// Try to fetch the user from the provider using existing tokens
	user, err := p.FetchUser(s)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
