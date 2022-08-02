package github

import (
	"errors"
	"net/http"

	"github.com/Lambels/autho"
	"github.com/google/go-github/v32/github"
)

var ErrNoGithubUser error = errors.New("unable to get github user")

func NewLoginHandler(cfg *autho.OAuth2Config) http.Handler {
	return autho.NewLoginHandler(cfg)
}

func NewTokenHandler(cfg *autho.OAuth2Config, errHandler, callbackHandler http.Handler) http.Handler {
	return autho.NewTokenHandler(cfg, errHandler, callbackHandler)
}

// NewCallbackHandler creates a new callback handler which on success sets the user under
// the request context and calls the terminal handler.
//
// The user type used by the default CallbackHandler is: https://github.com/google/go-github
func NewCallbackHandler(cfg *autho.OAuth2Config, errHandler, terminalHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		tkn, err := autho.TokenFromContext(r.Context())
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		// create a client and validate response.
		client := github.NewClient(
			cfg.OAuthConf.Client(r.Context(), tkn),
		)
		user, resp, err := client.Users.Get(r.Context(), "")
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}
		if resp.StatusCode != http.StatusOK {
			r = r.WithContext(autho.ContextWithError(r.Context(), ErrNoGithubUser))
			errHandler.ServeHTTP(w, r)
			return
		}
		if user == nil || user.ID == nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), ErrNoGithubUser))
			errHandler.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(autho.ContextWithUser(r.Context(), user))
		terminalHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
