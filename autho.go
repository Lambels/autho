package autho

import (
	"net/http"
	"time"
)

func NewApp(apps ...registerer) registerer {
	fn := func(mux multiplexer) {
		for _, app := range apps {
			app.Register(mux)
		}
	}

	return registererFunc(fn)
}

type CookieConfig struct {
	// Name sets the name of the cookie.
	Name string
	// Path Sets the path of the cookie which defaults to the path of the current responding
	// URL.
	Path string
	// Domain Sets the domain of the cookie which defaults to the domain of the current app.
	Domain string
	// Expires sets the expiry date of the cookie.
	Expires time.Time
	// MaxAge sets TTL for the cookie in seconds.
	MaxAge int
	// Secure indicates if the cookie will be sent through an HTTPS secure connection.
	Secure bool
	// HttpOnly indicates to the browser if the cookie is accessable by client-side scripts.
	HttpOnly bool
}

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
type app struct {
	CallbackURL string
	LoginURL    string

	// entrypoint for: callbackURL
	// responsible for validating state and obtaining token and passing token to callback handler.
	CallbackHandler http.Handler

	// entrypoint for: loginURL
	// responsible for setting state cookie and redirecting flow to oauth2 provider.
	LoginHandler http.Handler
}

func GetCookie(conf *CookieConfig, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(conf.Name)
	if err != nil {
		// if the cookie isnt present return a new one (not set).
		return newCookie(conf)
	}

	return cookie
}

func newCookie(conf *CookieConfig) *http.Cookie {
	return &http.Cookie{
		Name:     conf.Name,
		Path:     conf.Path,
		Domain:   conf.Domain,
		Expires:  conf.Expires,
		MaxAge:   conf.MaxAge,
		Secure:   conf.Secure,
		HttpOnly: conf.HttpOnly,
	}
}

// NewApp creates a new github oauth2 application.
func NewProviderApp(
	callbackURL, loginURL string,
	callbackHandler, loginHandler http.Handler,
) registerer {
	return &app{
		CallbackURL:     callbackURL,
		LoginURL:        loginURL,
		CallbackHandler: callbackHandler,
		LoginHandler:    loginHandler,
	}
}

func (a *app) Register(mux multiplexer) {
	// the url which the user hits and chains of the oauth2 flow.
	mux.Handle(a.LoginURL, a.LoginHandler)

	// the callback url.
	mux.Handle(a.CallbackURL, a.CallbackHandler)
}

type multiplexer interface {
	Handle(string, http.Handler)
}

type registerer interface {
	Register(multiplexer)
}

type registererFunc func(multiplexer)

func (f registererFunc) Register(mux multiplexer) {
	f(mux)
}
