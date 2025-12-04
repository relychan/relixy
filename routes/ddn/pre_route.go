// Package ddn implements restified endpoint handlers for DDN engine plugin.
package ddn

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/hasura/gotel"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/types"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const contentTypeHeader = "Content-Type"

// PreRoutePluginRequestBody represents the request body of the pre-route plugin request.
type PreRoutePluginRequestBody struct {
	Path   string          `json:"path"`
	Method string          `json:"method"`
	Query  string          `json:"query"`
	Body   json.RawMessage `json:"body"`
}

type preRoutePluginHandler struct {
	state *types.State
}

// NewPreRoutePluginHandler creates a pre-route plugin handler instance.
func NewPreRoutePluginHandler(state *types.State) *preRoutePluginHandler {
	return &preRoutePluginHandler{
		state: state,
	}
}

// ServeHTTP serves an HTTP request.
func (pr *preRoutePluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	logger := gotel.GetLogger(ctx)

	input, decoded := gohttps.DecodeRequestBody[PreRoutePluginRequestBody](w, r, span)
	if !decoded {
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

	resp, respBody, err := pr.state.ProxyClient.Execute(ctx, req) //nolint:bodyclose
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		logger.Error("failed to execute proxy request", slog.String("error", err.Error()))

		wErr := gohttps.WriteResponseError(w, err)
		if wErr != nil {
			gohttps.SetWriteResponseErrorAttribute(span, wErr)
		}

		return
	}

	respReader, ok := respBody.(io.ReadCloser)
	if ok {
		defer goutils.CatchWarnErrorFunc(respReader.Close)

		w.Header().Set(contentTypeHeader, resp.Header.Get(contentTypeHeader))
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, respReader)
	} else {
		err = gohttps.WriteResponseJSON(w, resp.StatusCode, respBody)
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		logger.Error("failed to write response", slog.String("error", err.Error()))
	}
}
