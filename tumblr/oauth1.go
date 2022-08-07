package tumblr

import (
	"net/http"

	"github.com/Lambels/autho"
	autho1 "github.com/Lambels/autho/oauth1"
	"github.com/dghubble/oauth1"
)

// NewCallbackHandler is a helper function that constructs
// a new callback handler using the default token handler tumblr.NewTokenHandler()
// wrapped arround the default tumblr.NewUserHandler().
//
// This method saves allot of boilerplate. For more customisable handlers construct your
// own callback handler by wrapping your own specific token handler around your own specific user handler.
func NewCallbackHandler(cfg *oauth1.Config, ckCfg *autho.CookieConfig, errHandler, terminalHandler http.Handler) http.Handler {
	return NewTokenHandler(
		cfg,
		ckCfg,
		errHandler,
		NewUserHandler(
			cfg,
			errHandler,
			terminalHandler,
		),
	)
}

// NewLoginHandler creates a new LoginHandler which is responsible for requesting the
// request token, if the provider requires to keep the request secret on the callback step
// the login handler is also responsible for setting the request secret in a cookie for later
// reading in the callback step. Afterwards the login handler is also responsible for redirecting
// the user to the provider.
//
// tumblr requires the request secret to be persisted throughout the callbacks.
func NewLoginHandler(cfg *oauth1.Config, ckCfg *autho.CookieConfig, errHandler http.Handler) http.Handler {
	return autho1.NewLoginHandler(cfg, ckCfg, errHandler)
}

// NewTokenHandler creates a new TokenHandler which is responsible for exchanging the
// request token, request secret (if ckCfg != nil) and verifier for the access token and
// access secret. Following the exchange the token handler is responsible to add the access token
// under the request ctx, calling on success the userHandler.
//
// tumblr requires to read the request secret in the login handler.
func NewTokenHandler(cfg *oauth1.Config, ckCfg *autho.CookieConfig, errHandler, userHandler http.Handler) http.Handler {
	return autho1.NewTokenHandler(cfg, ckCfg, errHandler, userHandler)
}

// NewUserHandler creates a new tumblr UserHandler resposnible for using the tokens provided
// by the TokenHandler in exchange for the users resource. The user resource is set under the
// request context.
//
//	user, ok := autho.UserFromContext(r.Context()).(*tumblr.User)
//
// The UserModel used by default by the tumblr.NewUserHandler is: https://www.tumblr.com/docs/en/api/v2#userinfo--get-a-users-information
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

		user, err := me(cfg.Client(
			r.Context(),
			tkn,
		))
		if err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}

		userCtx := autho.ContextWithUser(r.Context(), user)
		terminalHandler.ServeHTTP(w, r.WithContext(userCtx))
	}

	return http.HandlerFunc(f)
}
