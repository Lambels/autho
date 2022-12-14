package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/Lambels/autho"
	"golang.org/x/oauth2"
)

// NewLoginHandler creates a new LoginHandler which is resposible for setting a random
// value (state) to the state cookie. Afterwards the login handler is also
// responsible for redirecting the user to the provider for the users grant.
//
// LoginHandler -> Provider (obtain grant)
func NewLoginHandler(cfg *oauth2.Config, ckCfg *autho.CookieConfig) http.HandlerFunc {
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

		// redirect to provider url.
		redirectURL := cfg.AuthCodeURL(string(dst))
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

// NewTokenHandler creates a new TokenHandler which is the first handler in the chain responding
// to the callback from the provider, it is responsible for parsing the response for auth code
// and state then comparing the cookie state with the request state. Following the parsing the
// TokenHandler performs the token exchange and adds the token to the request context, calling on
// success the UserHandler.
//
// Provider -> TokenHandler -> UserHandler -> TermnialHandler
func NewTokenHandler(cfg *oauth2.Config, ckCfg *autho.CookieConfig, errHandler, userHandler http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = autho.DefaultFailureHandle
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		// parse auth code and state for token exchange.
		if err := r.ParseForm(); err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}
		state := r.Form.Get("state")
		authCode := r.Form.Get("code")

		if state == "" || authCode == "" {
			autho.PassError(errors.New("autho: auth code or state missing."), errHandler, w, r)
			return
		}

		// grab state cookie.
		ck, err := r.Cookie(ckCfg.Name)
		if err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}

		// validate any state mismatch.
		if state != string(ck.Value) {
			autho.PassError(errors.New("autho: request state and response state mismatch."), errHandler, w, r)
			return
		}

		// exchange auth code for token.
		tkn, err := cfg.Exchange(r.Context(), authCode)
		if err != nil {
			autho.PassError(err, errHandler, w, r)
			return
		}

		// redirect to user handler.
		tknCtx := ContextWithToken(r.Context(), tkn)
		userHandler.ServeHTTP(w, r.WithContext(tknCtx))
	}

	return http.HandlerFunc(fn)
}
