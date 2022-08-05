package facebook

import (
	"net/http"

	"github.com/Lambels/autho"
	autho2 "github.com/Lambels/autho/oauth2"
	fb "github.com/huandu/facebook/v2"
	"golang.org/x/oauth2"
)

// NewCallbackHandler is a helper function that constructs
// a new callback handler using the default token handler facebook.NewTokenHandler()
// wrapped arround the default facebook.NewUserHandler().
//
// This method saves allot of boilerplate. For more customisable handlers construct your
// own callback handler by wrapping your own specific token handler around your own specific user handler.
func NewCallbackHandler(cfg *oauth2.Config, ckCfg *autho.CookieConfig, errHandler, terminalHandler http.Handler) http.Handler {
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

// NewLoginHandler creates a new LoginHandler which is resposible for setting a random
// value (state) to the state cookie. Afterwards the login handler is also
// responsible for redirecting the user to the provider for the users grant.
func NewLoginHandler(ckCfg *autho.CookieConfig, oauthCfg *oauth2.Config) http.Handler {
	return autho2.NewLoginHandler(ckCfg, oauthCfg)
}

// NewTokenHandler creates a new TokenHandler which is the first handler in the chain responding
// to the callback from the provider, it is responsible for parsing the response for auth code
// and state then comparing the cookie state with the request state. Following the parsing the
// TokenHandler performs the token exchange and adds the token to the request context, calling on
// success the UserHandler.
func NewTokenHandler(cfg *oauth2.Config, ckCfg *autho.CookieConfig, errHandler, userHandler http.Handler) http.Handler {
	return autho2.NewTokenHandler(cfg, ckCfg, errHandler, userHandler)
}

// NewUserHandler creates a new facebook UserHandler resposnible for using the tokens provided
// by the TokenHandler in exchange for the users resource. The user resource is set under the
// request context.
//
//	user, ok := autho.UserFromContext(r.Context()).(*facebook.User)
//
// The UserModel used by default by the facebook.NewUserHandler is: https://developers.facebook.com/docs/graph-api/reference/user/#default-public-profile-fields
func NewUserHandler(cfg *oauth2.Config, errHandler, terminalHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		tkn, err := autho2.TokenFromContext(r.Context())
		if err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}

		httpClient := cfg.Client(r.Context(), tkn)
		session := &fb.Session{
			Version:    "v2.4",
			HttpClient: httpClient,
		}
		res, err := session.Get("/me", nil)
		if err != nil {
			// facebook api error.
			if e, ok := err.(*fb.Error); ok {
				autho.PassError(e, errHandler, w, r)
				return
			}

			// facebook unmarshal error.
			if e, ok := err.(*fb.UnmarshalError); ok {
				autho.PassError(e, errHandler, w, r)
				return
			}

			autho.PassError(autho.ErrNoUser, errHandler, w, r)
			return
		}

		var user User
		if err := res.Decode(&user); err != nil {
			autho.PassError(autho.ErrNoUser, errHandler, w, r)
			return
		}

		userCtx := autho.ContextWithUser(r.Context(), &user)
		terminalHandler.ServeHTTP(w, r.WithContext(userCtx))
	}

	return http.HandlerFunc(f)
}
