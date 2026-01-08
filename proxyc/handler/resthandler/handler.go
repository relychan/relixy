// Package resthandler evaluates and execute REST requests to the remote server.
package resthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/goutils"
	"github.com/relychan/goutils/httpheader"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
	"go.yaml.in/yaml/v4"
)

// RESTHandler implements the RelixyHandler interface for REST proxy.
type RESTHandler struct {
	method         string
	requestPath    string
	responseConfig proxyhandler.RelixyResponseConfig
	parameters     []*highv3.Parameter
}

// NewRESTHandler creates a RESTHandler from operation.
func NewRESTHandler(
	operation *highv3.Operation,
	rawProxyAction *yaml.Node,
	options *proxyhandler.NewRelixyHandlerOptions,
) (proxyhandler.RelixyHandler, error) {
	handler := &RESTHandler{
		method:     options.Method,
		parameters: openapi.MergeParameters(options.Parameters, operation.Parameters),
	}

	if rawProxyAction == nil {
		return handler, nil
	}

	var proxyAction RelixyRESTActionConfig

	err := rawProxyAction.Decode(&proxyAction)
	if err != nil {
		return nil, err
	}

	if proxyAction.Request != nil {
		if proxyAction.Request.Path != "" {
			handler.requestPath = proxyAction.Request.Path
		}
	}

	getEnvFunc := options.GetEnvFunc()

	responseConfig, err := proxyhandler.NewRelixyResponseConfig(
		proxyAction.Response,
		getEnvFunc,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize response config: %w", err)
	}

	handler.responseConfig = responseConfig

	return handler, nil
}

// Type returns type of the current handler.
func (*RESTHandler) Type() base_schema.RelixyActionType {
	return base_schema.ProxyTypeREST
}

// Handle resolves the HTTP request and proxies that request to the remote server.
func (re *RESTHandler) Handle(
	ctx context.Context,
	request *http.Request,
	options *proxyhandler.RelixyHandleOptions,
) (*http.Response, any, error) {
	requestPath := re.requestPath
	if requestPath == "" {
		requestPath = options.Path
	}

	if request.URL.RawQuery != "" {
		requestPath += "?" + request.URL.RawQuery
	}

	logAttrs := make([]slog.Attr, 0, 8)
	logAttrs = append(logAttrs, slog.String("path", requestPath))

	req := options.NewRequest(re.method, requestPath)

	if request.Body != nil {
		req.SetBody(request.Body)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		printDebugLog(
			ctx, request,
			"failed to execute request",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return resp, nil, err
	}

	logAttrs = append(logAttrs, slog.Int("response_status", resp.StatusCode))

	if re.responseConfig.IsZero() ||
		(resp.StatusCode < 200 && resp.StatusCode >= 300) ||
		resp.Header.Get(httpheader.ContentType) != httpheader.ContentTypeJSON {
		printDebugLog(
			ctx, request,
			resp.Status,
			logAttrs,
		)

		return resp, resp.Body, err
	}

	newResp, respBody, respLogAttrs, err := re.transformResponse(resp)
	logAttrs = append(logAttrs, respLogAttrs...)

	if err != nil {
		printDebugLog(
			ctx, request,
			"failed to transform response",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return resp, nil, err
	}

	printDebugLog(ctx, request, resp.Status, logAttrs)

	return newResp, respBody, err
}

func (re *RESTHandler) transformResponse( //nolint:revive
	resp *http.Response,
) (*http.Response, io.ReadCloser, []slog.Attr, error) {
	if re.responseConfig.Transform == nil || re.responseConfig.Transform.IsZero() {
		return resp, resp.Body, nil, nil
	}

	defer goutils.CatchWarnErrorFunc(resp.Body.Close)

	var responseBody any

	if resp.Body != nil {
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			return resp, nil, nil, fmt.Errorf("failed to decode http response: %w", err)
		}
	}

	responseLogAttrs := make([]slog.Attr, 0, 2)
	responseLogAttrs = append(responseLogAttrs, slog.Any("original_body", responseBody))

	transformedBody, err := re.responseConfig.Transform.Transform(responseBody)
	if err != nil {
		return nil, nil, nil, err
	}

	responseLogAttrs = append(responseLogAttrs, slog.Any("response_body", transformedBody))

	buf := new(bytes.Buffer)

	err = json.NewEncoder(buf).Encode(transformedBody)
	if err != nil {
		return resp, nil, nil, fmt.Errorf("failed to decode transformed response: %w", err)
	}

	return resp, io.NopCloser(buf), responseLogAttrs, err
}

func printDebugLog(
	ctx context.Context,
	request *http.Request,
	message string,
	attrs []slog.Attr,
) {
	logger := gotel.GetLogger(ctx)

	if !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	attrs = append(
		attrs,
		slog.String("type", "proxy-handler"),
		slog.String("handler_type", string(base_schema.ProxyTypeREST)),
		slog.String("request_url", request.URL.String()),
		otelutils.NewHeaderLogGroupAttrs(
			"request_headers",
			otelutils.NewTelemetryHeaders(request.Header),
		),
	)

	logger.LogAttrs(ctx, slog.LevelDebug, message, attrs...)
}
