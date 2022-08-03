package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/Lambels/autho"
	"golang.org/x/oauth2"
)

func NewLoginHandler(ckCfg *autho.CookieConfig, oauthCfg *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get any existing cookie or create a new one.
		ck := autho.GetCookie(ckCfg, r)

		// generate random state.
		buf := make([]byte, 32)
		rand.Read(buf)
		dst := make([]byte, 32)
		base64.RawURLEncoding.Encode(dst, buf)
		ck.Value = string(dst)

		// set state.
		http.SetCookie(w, ck)
		r = r.WithContext(ContextWithState(r.Context(), dst))

		// redirect to provider url.
		redirectURL := oauthCfg.AuthCodeURL(string(dst))
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func NewTokenHandler(cfg *oauth2.Config, errHandler, userHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		// parse auth code and state for token exchange.
		if err := r.ParseForm(); err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}
		state := r.Form.Get("state")
		authCode := r.Form.Get("code")

		if state == "" || authCode == "" {
			r = r.WithContext(autho.ContextWithError(r.Context(), errors.New("auth code or state missing.")))
			errHandler.ServeHTTP(w, r)
			return
		}

		// grab request state.
		requestState, err := StateFromContext(r.Context())
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		// validate any state mismatch.
		if state != string(requestState) {
			r = r.WithContext(autho.ContextWithError(r.Context(), errors.New("request state and response state mismatch.")))
			errHandler.ServeHTTP(w, r)
			return
		}

		// exchange auth code for token.
		token, err := cfg.Exchange(r.Context(), authCode)
		if err != nil {
			r = r.WithContext(autho.ContextWithError(r.Context(), err))
			errHandler.ServeHTTP(w, r)
			return
		}

		// redirect to user handler.
		r = r.WithContext(ContextWithToken(r.Context(), token))
		userHandler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
