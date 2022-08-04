package autho

import (
	"context"
	"errors"
)

type errKey struct{}

// ContextWithError adds err to the context to be used by the err handler.
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

// ContextWithUser adds the user to the context to be used by the terminal handler.
func ContextWithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// UserFromContext harvests the user struct from the request context.
// When the struct is returned if no error is present
// parse it to the specific user you are expecting from the user handler (check your providers
// user handler for extra info)
func UserFromContext(ctx context.Context) (interface{}, error) {
	user := ctx.Value(userKey{})
	if user == nil {
		return nil, errors.New("autho: user parameter not ser")
	}

	return user, nil
}
