package oauth1

import (
	"context"
	"errors"

	"github.com/dghubble/oauth1"
)

type requestSecretKey struct{}

func ContextWithRequestSecret(ctx context.Context, secret string) context.Context {
	return context.WithValue(ctx, requestSecretKey{}, secret)
}

func RequestSecretFromContext(ctx context.Context) (secret string, err error) {
	secret, ok := ctx.Value(requestSecretKey{}).(string)
	if !ok {
		return "", errors.New("request secret not set")
	}
	return secret, nil
}

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
