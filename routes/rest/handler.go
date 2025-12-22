// Package rest implements restified endpoint handlers for DDN engine plugin.
package rest

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/hasura/gotel"
	"github.com/relychan/gohttps/httputils"
	"github.com/relychan/goutils"
	"github.com/relychan/goutils/httpheader"
	"github.com/relychan/relyx/types"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type restHandler struct {
	state *types.State
}

// NewRESTHandler creates a REST handler instance.
func NewRESTHandler(state *types.State) *restHandler {
	return &restHandler{
		state: state,
	}
}

// ServeHTTP serves an HTTP request.
func (rh *restHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	logger := gotel.GetLogger(ctx)

	resp, respBody, err := rh.state.ProxyClient.Execute(ctx, r) //nolint:bodyclose
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		logger.Error("failed to execute proxy request", slog.String("error", err.Error()))

		wErr := httputils.WriteResponseError(w, err)
		if wErr != nil {
			httputils.SetWriteResponseErrorAttribute(span, wErr)
		}

		return
	}

	respReader, ok := respBody.(io.ReadCloser)
	if ok {
		defer goutils.CatchWarnErrorFunc(respReader.Close)

		w.Header().Set(httpheader.ContentType, resp.Header.Get(httpheader.ContentType))
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, respReader)
	} else {
		err = httputils.WriteResponseJSON(w, resp.StatusCode, respBody)
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		logger.Error("failed to write response", slog.String("error", err.Error()))
	}
}
