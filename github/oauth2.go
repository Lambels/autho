package github

import (
	"net/http"

	"github.com/Lambels/autho"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// NewDefaultCallbackHandler is a helper function that constructs
// a new callback handler using the default token handler github.NewTokenHandler()
// wrapped arround the default github.NewUserHandler().
//
// This method saves allot of boilerplate. For more customisable handlers construct your
// own callback handler by wrapping it around your own specific token handler.
func NewDefaultCallbackHandler(oauthCfg *oauth2.Config, errHandler, terminalHandler http.Handler) http.Handler {
	return NewTokenHandler(
		oauthCfg,
		errHandler,
		NewUserHandler(
			oauthCfg,
			errHandler,
			terminalHandler,
		),
	)
}

// NewLoginHandler creates a new LoginHandler which is resposible for setting a random
// value (state) to the request context and state cookie. Afterwards the login handler is also
// responsible for redirecting the user to the provider for the users grant.
func NewLoginHandler(ckCfg *autho.CookieConfig, oauthCfg *oauth2.Config) http.Handler {
	return autho.NewLoginHandler(ckCfg, oauthCfg)
}

// NewTokenHandler creates a new TokenHandler which is the first handler in the chain responding
// to the callback from the provider, it is responsible for parsing the response for auth code
// and state then comparing the ctx state with the request state. Following the parsing the
// TokenHandler performs the token exchange and adds the token to the request context, calling on
// success the UserHandler.
func NewTokenHandler(cfg *oauth2.Config, errHandler, callbackHandler http.Handler) http.Handler {
	return autho.NewTokenHandler(cfg, errHandler, callbackHandler)
}

// NewUserHandler creates a new github UserHandler resposnible for using the tokens provided
// by the TokenHandler in exchange for the users resource. The user resource is set under the
// request context.
//
//	user, ok := autho.UserFromContext(r.Context()).(*github.User)
//
// The UserModel used by default by the github.NewUserHandler is: https://pkg.go.dev/github.com/google/go-github/v45/github#User
func NewUserHandler(cfg *oauth2.Config, errHandler, terminalHandler http.Handler) http.Handler {
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
			cfg.Client(r.Context(), tkn),
		)
		user, resp, err := client.Users.Get(r.Context(), "")
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}
		if resp.StatusCode != http.StatusOK {
			r = r.WithContext(autho.ContextWithError(r.Context(), autho.ErrNoUser))
			errHandler.ServeHTTP(w, r)
			return
		}
		if user == nil || user.ID == nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), autho.ErrNoUser))
			errHandler.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(autho.ContextWithUser(r.Context(), user))
		terminalHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
