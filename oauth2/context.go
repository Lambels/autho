package oauth2

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
)

type stateKey struct{}

func ContextWithState(ctx context.Context, buf []byte) context.Context {
	return context.WithValue(ctx, stateKey{}, buf)
}

func StateFromContext(ctx context.Context) ([]byte, error) {
	state, ok := ctx.Value(stateKey{}).([]byte)
	if !ok {
		return []byte{}, errors.New("state parameter not set")
	}

	return state, nil
}

type tokenKey struct{}

func ContextWithToken(ctx context.Context, tkn *oauth2.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, tkn)
}

func TokenFromContext(ctx context.Context) (*oauth2.Token, error) {
	tkn, ok := ctx.Value(tokenKey{}).(*oauth2.Token)
	if !ok {
		return &oauth2.Token{}, errors.New("token parameter not set")
	}

	return tkn, nil
}
