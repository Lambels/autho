package oauth1

import (
	"net/http"

	"github.com/Lambels/autho"
	"github.com/dghubble/oauth1"
)

func NewLoginHandler(cfg *oauth1.Config, ckCfg *autho.CookieConfig, errHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		reqToken, reqSecret, err := cfg.RequestToken()
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		if ckCfg != nil {
			ck := autho.GetCookie(ckCfg, r)
			ck.Value = reqSecret
			http.SetCookie(w, ck)
		}

		authURL, err := cfg.AuthorizationURL(reqToken)
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, authURL.String(), http.StatusFound)
	}

	return http.HandlerFunc(f)
}

func NewTokenHandler(cfg *oauth1.Config, ckCfg *autho.CookieConfig, errHandler, userHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		reqToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		var reqSecret string
		if ckCfg != nil {
			ck, ckErr := r.Cookie(ckCfg.Name)
			reqSecret, err = ck.Value, ckErr
		} else {
			reqSecret, err = "", nil
		}
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		accessToken, accessSecret, err := cfg.AccessToken(reqToken, reqSecret, verifier)
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		tkn := oauth1.NewToken(accessToken, accessSecret)
		r = r.WithContext(RequestWithToken(r.Context(), tkn))
		userHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
