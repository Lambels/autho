package oauth

import (
	"net/http"

	"golang.org/x/oauth2"
)

func NewLoginHandler(cfg *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func NewTokenHandler(cfg *oauth2.Config, errHandler, next http.Handler) http.Handler {
	if errHandler == nil {
		errHandler = DefaultFailureHandle
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
