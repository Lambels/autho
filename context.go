package autho

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
)

type errKey struct{}

func ContextWithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errKey{}, err)
}

// ErrorFromContext returns any error from ctx.
// If an error is present then its parsed and returned else nil is returned.
func ErrorFromContext(ctx context.Context) error {
	val := ctx.Value(errKey{})
	if val != nil {
		// safe data type conversion since values under the errKey are locked under ContextWithError.
		return val.(error)
	}
	return nil
}

type userKey struct{}

func ContextWithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// UserFromContext retrieves the user struct under the user key.
// When the struct is returned parse it to the specific user you are expecting, if the parcing
// fails no user is set.
func UserFromContext(ctx context.Context) interface{} {
	return ctx.Value(userKey{})
}

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
