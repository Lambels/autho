package oauth

import (
	"net/http"
)

// DefaultFailureHandle sends a response with error code: 400 (Bad Request) and the error text.
var DefaultFailureHandle http.HandlerFunc = failureHandler

func failureHandler(w http.ResponseWriter, r *http.Request) {
	if err := ErrorFromContext(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

// App provides an abstraction from oauth 1 or 2 for each oauth aplication.
type App struct {
	CallbackURL string
	LoginURL    string

	// entrypoint for: callbackURL
	// responsible for validating state and obtaining token and passing token to callback handler.
	CallbackHandler http.Handler

	// entrypoint for: loginURL
	// responsible for setting state cookie and redirecting flow to oauth2 provider.
	LoginHandler http.Handler
}

// NewApp creates a new github oauth2 application.
func NewApp(
	callbackURL, loginURL string,
	callbackHandler, loginHandler http.Handler,
) Registerer {
	return &App{
		CallbackURL:     callbackURL,
		LoginURL:        loginURL,
		CallbackHandler: callbackHandler,
		LoginHandler:    loginHandler,
	}
}

func (a *App) Register(mux Multiplexer) {
	// the url which the user hits and chains of the oauth2 flow.
	mux.Handle(a.LoginURL, a.LoginHandler)

	// the callback url.
	mux.Handle(a.CallbackURL, a.CallbackHandler)
}

type Multiplexer interface {
	Handle(string, http.Handler)
}

type Registerer interface {
	Register(Multiplexer)
}

type RegistererFunc func(Multiplexer)

func (f RegistererFunc) Register(mux Multiplexer) {
	f(mux)
}
