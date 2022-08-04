package autho

import (
	"errors"
	"net/http"
	"time"
)

// ErrNoUser represents the user handler not being able to reach the user resoursce.
var ErrNoUser error = errors.New("autho: unable to get user from provider")

// DefaultFailureHandle sends a response with error code: 400 (Bad Request) and the error text.
var DefaultFailureHandle http.HandlerFunc = failureHandler

// NewApp creates a new autho app which consists of multiple providers (OAuth 1 or 2),
// the new app is used to register the providers' callback and login URL to the provided
// multiplexer.
//
//	mux := http.NewServeMux()
//	app := autho.NewApp()
//	app.Register(mux)
func NewApp(apps ...Registerer) Registerer {
	fn := func(mux multiplexer) {
		for _, app := range apps {
			app.Register(mux)
		}
	}

	return registererFunc(fn)
}

// app provides an abstraction from oauth 1 or 2 for each oauth aplication.
type app struct {
	callbackURL string
	loginURL    string

	// entrypoint for: callbackURL
	callbackHandler http.Handler

	// entrypoint for: loginURL
	loginHandler http.Handler
}

// Register registers the callback and login URL to mux.
func (a *app) Register(mux multiplexer) {
	// the url which the user hits and chains of the oauth flow.
	mux.Handle(a.loginURL, a.loginHandler)

	// the callback url.
	mux.Handle(a.callbackURL, a.callbackHandler)
}

// NewProviderApp creates a new provider application.
func NewProviderApp(
	callbackURL, loginURL string,
	callbackHandler, loginHandler http.Handler,
) Registerer {
	return &app{
		callbackURL:     callbackURL,
		loginURL:        loginURL,
		callbackHandler: callbackHandler,
		loginHandler:    loginHandler,
	}
}

// CookieConfig represents a config used by handlers to create cookies throughout the
// oauth flow.
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

// GetCookie gets a cookie from the request if possible or returns a new cookie
// structured after the config.
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

func failureHandler(w http.ResponseWriter, r *http.Request) {
	if err := ErrorFromContext(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

type multiplexer interface {
	Handle(string, http.Handler)
}

type Registerer interface {
	Register(multiplexer)
}

type registererFunc func(multiplexer)

func (f registererFunc) Register(mux multiplexer) {
	f(mux)
}
