package authn

import (
	"context"

	"github.com/relychan/goutils"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct{}

func (contextKey) String() string {
	return "auth context key"
}

var authContextKey = &contextKey{}

// NewAuthContext creates a new context with an authenticated session.
func NewAuthContext(parentContext context.Context, value any) context.Context {
	return context.WithValue(parentContext, authContextKey, value)
}

// GetAuthContext gets the authenticated session from context.
func GetAuthContext[T any](ctx context.Context) (T, error) { //nolint:ireturn
	rawValue := ctx.Value(authContextKey)

	if rawValue == nil {
		var zeroValue T

		return zeroValue, goutils.NewUnauthorizedError()
	}

	value, ok := rawValue.(T)
	if !ok {
		return value, goutils.NewServerError(goutils.ErrorDetail{
			Detail: "unexpected authenticated context type",
		})
	}

	return value, nil
}
