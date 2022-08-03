package twitter

import (
	"net/http"

	"github.com/Lambels/autho"
	autho1 "github.com/Lambels/autho/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

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

func NewLoginHandler(cfg *oauth1.Config, errHandler http.Handler) http.Handler {
	return autho1.NewLoginHandler(cfg, nil, errHandler)
}

func NewTokenHandler(cfg *oauth1.Config, errHandler, userHandler http.Handler) http.Handler {
	return autho1.NewTokenHandler(cfg, nil, errHandler, userHandler)
}

func NewUserHandler(cfg *oauth1.Config, errHandler, terminalHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		tkn, err := autho1.TokenFromContext(r.Context())
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
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
			r = r.WithContext(autho.ContextWithError(r.Context(), autho.ErrNoUser))
			errHandler.ServeHTTP(w, r)
			return
		}
		if user == nil || user.ID == 0 {
			r = r.WithContext(autho.ContextWithError(r.Context(), autho.ErrNoUser))
			errHandler.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(autho.ContextWithUser(r.Context(), user))
		terminalHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
