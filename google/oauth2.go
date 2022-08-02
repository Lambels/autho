package google

import (
	"net/http"

	"github.com/Lambels/autho"
	"golang.org/x/oauth2"
	googleOauth "google.golang.org/api/oauth2/v2"
)

func NewLoginHandler(ckCfg *autho.CookieConfig, oauthCfg *oauth2.Config) http.Handler {
	return autho.NewLoginHandler(ckCfg, oauthCfg)
}

func NewTokenHandler(cfg *oauth2.Config, errHandler, callbackHandler http.Handler) http.Handler {
	return autho.NewTokenHandler(cfg, errHandler, callbackHandler)
}

func NewCallbackHandler(cfg *oauth2.Config, errHandler, terminalHandler http.Handler) http.Handler {
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

		service, err := googleOauth.New(
			cfg.Client(r.Context(), tkn),
		)
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		userInfo, err := service.Userinfo.Get().Do()
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}
		if userInfo.Id == "" {
			r = r.WithContext(autho.ContextWithError(r.Context(), autho.ErrNoUser))
			errHandler.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(autho.ContextWithUser(r.Context(), userInfo))
		terminalHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
