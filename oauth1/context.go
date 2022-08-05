package oauth1

import (
	"context"
	"errors"

	"github.com/dghubble/oauth1"
)

type tokenKey struct{}

// ContextWithToken is used by the token handler (default: oauth1.NewTokenHandler()) to
// set the oauth1 tokens under the context to be used by the user handler in exchange for
// the user resource.
func ContextWithToken(ctx context.Context, tkn *oauth1.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, tkn)
}

// TokenFromContext is used to harvest the tokens from the request context by the
// user handler and preform the token/user-resource exchange.
func TokenFromContext(ctx context.Context) (*oauth1.Token, error) {
	tkn, ok := ctx.Value(tokenKey{}).(*oauth1.Token)
	if !ok {
		return &oauth1.Token{}, errors.New("autho: token parameter not set")
	}

	return tkn, nil
}
