// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
