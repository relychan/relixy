package proxyc

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Execute routes and proxies the request to the remote server.
func (pc *ProxyClient) Execute(
	ctx context.Context,
	req *http.Request,
) (*http.Response, any, error) {
	ctx, span := tracer.Start(ctx, "Proxy")
	defer span.End()

	span.SetAttributes(
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.URLOriginal(req.URL.String()),
	)

	requestPath := req.URL.Path

	if pc.clientOptions.BasePath != "" && req.URL.Path != "" {
		// The URL path may omit the slash character
		if req.URL.Path[0] == '/' {
			requestPath = strings.TrimPrefix(req.URL.Path, pc.clientOptions.BasePath)
		} else {
			requestPath = strings.TrimPrefix(req.URL.Path, pc.clientOptions.BasePath[1:])
		}
	}

	span.SetAttributes(semconv.URLPath(requestPath))

	route := pc.node.FindRoute(requestPath, req.Method)
	if route == nil {
		span.SetStatus(codes.Error, "request path or method does not exist")

		return nil, nil, goutils.RFC9457Error{
			Status:   http.StatusNotFound,
			Title:    "Resource Not Found",
			Instance: req.URL.Path,
		}
	}

	span.SetAttributes(attribute.String("http.request.proxy.type", string(route.Handler.Type())))

	options := &schema.RelyProxyHandleOptions{
		Settings:       &pc.metadata.Settings,
		ParamValues:    route.ParamValues,
		HTTPClient:     pc.lbClient,
		DefaultHeaders: pc.defaultHeaders,
		Path:           requestPath,
	}

	response, responseBody, err := route.Handler.Handle(ctx, req, options)
	if err != nil {
		span.SetStatus(codes.Error, "proxy failed")
		span.RecordError(err)

		var rfc9457Error goutils.RFC9457Error

		if errors.As(err, &rfc9457Error) {
			rfc9457Error.Instance = req.URL.Path
		} else {
			rfc9457Error = goutils.RFC9457Error{
				Status:   http.StatusInternalServerError,
				Title:    http.StatusText(http.StatusInternalServerError),
				Detail:   err.Error(),
				Instance: req.URL.Path,
			}
		}

		return nil, nil, rfc9457Error
	}

	span.SetStatus(codes.Ok, "")

	return response, responseBody, nil
}
