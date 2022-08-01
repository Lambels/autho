package oauth

import (
	"net/http"

	"github.com/Lambels/autho/internal"
	"golang.org/x/oauth2"
)

type app struct {
	callbackURL string
	loginURL    string
	// the oauth2 configuration.
	config *oauth2.Config

	// entrypoint for: callbackURL
	// -> callbackHandler / failiureHandler
	// responsible for validating state and obtaining token and passing token to callback handler.
	tokenHandler internal.Middleware
	// responsible for token exchange.
	// -> successHandler / failiureHandler
	callbackHandler internal.Middleware

	// entrypoint for: loginURL
	// responsible for setting state cookie and redirecting flow to oauth2 provider.
	loginHandler http.HandlerFunc
}

func (a *app) register(mux multiplexer) {
	// the url which the user hits and chains of the oauth2 flow.
	mux.HandleFunc(a.loginURL, a.loginHandler)

	// the callback url.
	mux.HandleFunc(a.callbackURL)
}

// loginHandler is the default oauth2 handler for the login URL.
func loginHandler(w http.ResponseWriter, r *http.Request) {

}

func failiureHandler(w http.ResponseWriter, r *http.Request) {

}
