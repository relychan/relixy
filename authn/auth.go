// Package authn defines the middleware and functions to authenticate HTTP requests.
package authn

import (
	"net/http"

	"github.com/go-viper/mapstructure/v2"
	"github.com/relychan/gohttps/httputils"
	"github.com/relychan/goutils"
	"github.com/relychan/rely-auth/auth"
	"github.com/relychan/rely-auth/auth/authmode"
	"go.opentelemetry.io/otel/trace"
)

// AuthMiddleware authenticates requests and extracts sessions variables from the authenticated credentials.
func AuthMiddleware[T any](authManager *auth.RelyAuthManager) func(http.Handler) http.Handler {
	if authManager == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return &authMiddleware[T]{
			manager: authManager,
			next:    next,
		}
	}
}

type authMiddleware[T any] struct {
	manager *auth.RelyAuthManager
	next    http.Handler
}

// ServeHTTP implements the http.Handler interface to authenticate the request.
func (am *authMiddleware[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authRequest := authmode.AuthenticateRequestData{
		URL:     r.URL.String(),
		Headers: goutils.ExtractHeaders(r.Header),
	}

	ctx := r.Context()

	sessionVariables, err := am.manager.Authenticate(r.Context(), authRequest)
	if err != nil {
		writeErr := httputils.WriteResponseError(w, err)
		if writeErr != nil {
			httputils.SetWriteResponseErrorAttribute(trace.SpanFromContext(r.Context()), writeErr)
		}

		return
	}

	var result T

	_, ok := any(result).(map[string]any)
	if ok {
		newReq := r.WithContext(NewAuthContext(ctx, sessionVariables))
		am.next.ServeHTTP(w, newReq)

		return
	}

	// use mapstructure to map session variables to the custom structure.
	err = mapstructure.Decode(sessionVariables, &result)
	if err != nil {
		respBody := goutils.NewServerError()
		respBody.Detail = err.Error()
		respBody.Instance = r.URL.Path

		writeErr := httputils.WriteResponseJSON(w, respBody.Status, respBody)
		if writeErr != nil {
			httputils.SetWriteResponseErrorAttribute(trace.SpanFromContext(r.Context()), writeErr)
		}

		return
	}

	newReq := r.WithContext(NewAuthContext(ctx, result))
	am.next.ServeHTTP(w, newReq)
}
