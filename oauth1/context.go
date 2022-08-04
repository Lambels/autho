package oauth1

import (
	"context"
	"errors"

	"github.com/dghubble/oauth1"
)

type tokenKey struct{}

func RequestWithToken(ctx context.Context, token *oauth1.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

func TokenFromContext(ctx context.Context) (*oauth1.Token, error) {
	tkn, ok := ctx.Value(tokenKey{}).(*oauth1.Token)
	if !ok {
		return nil, errors.New("token not ser")
	}

	return tkn, nil
}
