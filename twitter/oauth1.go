package twitter

import (
	"net/http"

	"github.com/Lambels/autho"
	autho1 "github.com/Lambels/autho/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// NewCallbackHandler is a helper function that constructs
// a new callback handler using the default token handler twitter.NewTokenHandler()
// wrapped arround the default twitter.NewUserHandler().
//
// This method saves allot of boilerplate. For more customisable handlers construct your
// own callback handler by wrapping your own specific token handler around your own specific user handler.
func NewCallbackHandler(cfg *oauth1.Config, errHandler, terminalHandler http.Handler) http.Handler {
	return NewTokenHandler(
		cfg,
		errHandler,
		NewUserHandler(
			cfg,
			errHandler,
			terminalHandler,
		),
	)
}

// NewLoginHandler creates a new LoginHandler which is responsible for requesting the
// request token, twitter doesent need the request secret to be persisted to the callback step.
// Afterwards the login handler is also responsible for redirecting the user to the provider.
func NewLoginHandler(cfg *oauth1.Config, errHandler http.Handler) http.Handler {
	return autho1.NewLoginHandler(cfg, nil, errHandler)
}

// NewTokenHandler creates a new TokenHandler which is responsible for exchanging the
// request token and verifier for the access token and access secret.
func NewTokenHandler(cfg *oauth1.Config, errHandler, userHandler http.Handler) http.Handler {
	return autho1.NewTokenHandler(cfg, nil, errHandler, userHandler)
}

// NewUserHandler creates a new twitter UserHandler resposnible for using the tokens provided
// by the TokenHandler in exchange for the users resource. The user resource is set under the
// request context.
//
//	user, ok := autho.UserFromContext(r.Context()).(*twitter.User)
//
// The UserModel used by default by the twitter.NewUserHandler is: https://pkg.go.dev/github.com/dghubble/go-twitter/twitter#User
func NewUserHandler(cfg *oauth1.Config, errHandler, terminalHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		tkn, err := autho1.TokenFromContext(r.Context())
		if err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}

		httpClient := cfg.Client(r.Context(), tkn)
		client := twitter.NewClient(httpClient)

		user, resp, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
			IncludeEntities: twitter.Bool(false),
			SkipStatus:      twitter.Bool(true),
			IncludeEmail:    twitter.Bool(false),
		})
		if err != nil || resp.StatusCode != http.StatusOK {
			autho.PassError(autho.ErrNoUser, errHandler, w, r)
			return
		}
		if user == nil || user.ID == 0 {
			autho.PassError(autho.ErrNoUser, errHandler, w, r)
			return
		}

		r = r.WithContext(autho.ContextWithUser(r.Context(), user))
		terminalHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
