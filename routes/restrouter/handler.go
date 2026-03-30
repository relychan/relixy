// Package restrouter implements restified endpoint handlers for DDN engine plugin.
package restrouter

import (
	"log/slog"
	"net/http"

	"github.com/hasura/gotel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type restHandler struct {
	state *State
}

// NewRESTHandler creates a REST handler instance.
func NewRESTHandler(state *State) *restHandler {
	return &restHandler{
		state: state,
	}
}

// ServeHTTP serves an HTTP request.
func (rh *restHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	logger := gotel.GetLogger(ctx)

	_, err := rh.state.ProxyClient.Stream(w, r) //nolint:bodyclose
	if err == nil {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)

	logger.Error("failed to write response", slog.String("error", err.Error()))
}
