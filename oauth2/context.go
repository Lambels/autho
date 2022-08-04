package oauth2

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
)

type tokenKey struct{}

// ContextWithToken is used by the token handler (default: oauth2.NewTokenHandler()) to
// set the oauth2 tokens under the context to be used by the user handler in exchange for
// the user resource.
func ContextWithToken(ctx context.Context, tkn *oauth2.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, tkn)
}

// TokenFromContext is used to harvest the tokens from the request context by the
// user handler and preform the token/user-resource exchange.
func TokenFromContext(ctx context.Context) (*oauth2.Token, error) {
	tkn, ok := ctx.Value(tokenKey{}).(*oauth2.Token)
	if !ok {
		return &oauth2.Token{}, errors.New("autho: token parameter not set")
	}

	return tkn, nil
}
