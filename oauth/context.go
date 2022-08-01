package oauth

import (
	"context"
	"errors"
)

type errKey struct{}

func ContextWithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errKey{}, err)
}

func ErrorFromContext(ctx context.Context) error {
	val := ctx.Value(errKey{})
	if val != nil {
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
