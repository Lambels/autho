package oauth1

import (
	"net/http"

	"github.com/Lambels/autho"
	"github.com/dghubble/oauth1"
)

// NewLoginHandler creates a new LoginHandler which is responsible for requesting the
// request token, if the provider requires to keep the request secret on the callback step
// the login handler is also responsible for setting the request secret in a cookie for later
// reading in the callback step. Afterwards the login handler is also responsible for redirecting
// the user to the provider.
//
// LoginHandler -> Provider (obtain grant)
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

		// if a cookie config is provided, it flags that the provider needs the req secret
		// in the callback step, add it to a cookie.
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

// NewTokenHandler creates a new TokenHandler which is responsible for exchanging the
// request token, request secret (if ckCfg != nil) and verifier for the access token and
// access secret. Following the exchange the token handler is responsible to add the access token
// under the request ctx, calling on success the userHandler.
//
// If ckCfg is provider (!= nil) then the provider needs also the request token for the
// exchange. The upstream handler is required (in the redirect phase) to add the request
// secret to the cookie. Read the request secret from the cookie and pass it to the exchange
// else pass an empty string to the exchange.
//
// Provider -> TokenHandler -> UserHandler -> TerminalHandler
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

		// set request secret if ckCfg isnt nill.
		var reqSecret string
		if ckCfg != nil {
			ck, err := r.Cookie(ckCfg.Name)
			if err != nil {
				r = r.WithContext(autho.ContextWithError(r.Context(), err))
				errHandler.ServeHTTP(w, r)
				return
			}
			reqSecret = ck.Value
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
