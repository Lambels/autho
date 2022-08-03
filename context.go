package autho

import (
	"context"
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
