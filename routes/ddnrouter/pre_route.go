// Package ddnrouter implements restified endpoint handlers for DDN engine plugin.
package ddnrouter

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/hasura/gotel"
	"github.com/relychan/gohttps/httputils"
	"github.com/relychan/goutils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// PreRoutePluginRequestBody represents the request body of the pre-route plugin request.
type PreRoutePluginRequestBody struct {
	Path   string          `json:"path"`
	Method string          `json:"method"`
	Query  string          `json:"query"`
	Body   json.RawMessage `json:"body"`
}

type preRoutePluginHandler struct {
	state *State
}

// NewPreRoutePluginHandler creates a pre-route plugin handler instance.
func NewPreRoutePluginHandler(state *State) *preRoutePluginHandler {
	return &preRoutePluginHandler{
		state: state,
	}
}

// ServeHTTP serves an HTTP request.
func (pr *preRoutePluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	logger := gotel.GetLogger(ctx)

	input, decoded := httputils.DecodeRequestBody[PreRoutePluginRequestBody](w, r, span)
	if !decoded {
		return
	}

	proxyClient := pr.state.FindProxyClient(input.Path)
	if proxyClient == nil {
		err := goutils.NewNotFoundError()

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		logger.Warn("failed to execute proxy request", slog.String("error", err.Error()))

		wErr := httputils.WriteResponseError(w, err)
		if wErr != nil {
			httputils.SetWriteResponseErrorAttribute(span, wErr)
		}

		return
	}

	req := &http.Request{
		Method: input.Method,
		URL: &url.URL{
			Path:     input.Path,
			RawQuery: input.Query,
		},
	}

	if len(input.Body) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(input.Body))
	}

	_, err := proxyClient.Stream(w, req) //nolint:bodyclose
	if err == nil {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)

	logger.Error("failed to write response", slog.String("error", err.Error()))
}
