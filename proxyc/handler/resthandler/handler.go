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
	"github.com/relychan/gohttpc"
	"github.com/relychan/goutils"
	"github.com/relychan/goutils/httpheader"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/schema/openapi"
	"go.yaml.in/yaml/v4"
)

// RESTHandler implements the RelixyHandler interface for REST proxy.
type RESTHandler struct {
	method         string
	customRequest  *customRESTRequest
	customResponse *customRESTResponse
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

	getEnvFunc := options.GetEnvFunc()

	handler.customRequest, err = newCustomRESTRequestFromConfig(proxyAction.Request, getEnvFunc)
	if err != nil {
		return nil, err
	}

	handler.customResponse, err = newCustomRESTResponse(
		proxyAction.Response,
		getEnvFunc,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize response config: %w", err)
	}

	return handler, nil
}

// Type returns type of the current handler.
func (*RESTHandler) Type() proxyhandler.ProxyActionType {
	return ProxyActionTypeREST
}

// Handle resolves the HTTP request and proxies that request to the remote server.
func (re *RESTHandler) Handle(
	ctx context.Context,
	request *http.Request,
	options *proxyhandler.RelixyHandleOptions,
) (*http.Response, any, error) {
	req, logAttrs, err := re.transformRequest(request, options)
	if err != nil {
		printDebugLog(
			ctx, request,
			"failed to evaluate request",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return nil, nil, err
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

	if re.customResponse == nil || re.customResponse.IsZero() ||
		(resp.StatusCode < 200 || resp.StatusCode >= 300) ||
		resp.Header.Get(httpheader.ContentType) != httpheader.ContentTypeJSON {
		printDebugLog(
			ctx, request,
			resp.Status,
			logAttrs,
		)

		return resp, resp.Body, err
	}

	newResp, respLogAttrs, err := re.transformResponse(resp)
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

	return newResp, newResp.Body, err
}

func (re *RESTHandler) transformRequest(
	request *http.Request,
	options *proxyhandler.RelixyHandleOptions,
) (*gohttpc.RequestWithClient, []slog.Attr, error) {
	requestPath := options.Path

	if re.customRequest != nil && re.customRequest.Path != "" {
		requestPath = re.customRequest.Path
	}

	if request.URL.RawQuery != "" {
		requestPath += "?" + request.URL.RawQuery
	}

	logAttrs := make([]slog.Attr, 0, 8)
	logAttrs = append(logAttrs, slog.String("path", requestPath))

	req := options.NewRequest(re.method, requestPath)

	if re.customRequest == nil || re.customRequest.IsZero() {
		// Proxies the raw request to the remote service when there is no request.
		if request.Body != nil && request.Body != http.NoBody {
			req.SetBody(request.Body)
		}

		return req, logAttrs, nil
	}

	requestData, alreadyRead, err := proxyhandler.NewRequestTemplateData(
		request,
		request.Header.Get(httpheader.ContentType),
		options.ParamValues,
	)
	if err != nil {
		return nil, logAttrs, err
	}

	rawRequestData := requestData.ToMap()

	for key, headerEntry := range re.customRequest.Headers {
		value, err := headerEntry.EvaluateString(rawRequestData)
		if err != nil {
			return nil, logAttrs, fmt.Errorf("failed to transform request header %s: %w", key, err)
		}

		if value != nil && *value != "" {
			req.Header().Set(key, *value)
		}
	}

	if !alreadyRead || re.customRequest.Body == nil || re.customRequest.Body.IsZero() {
		if !alreadyRead {
			// unsupported content types will be ignored,
			// the client proxies the raw request to the remote service.
			req.SetBody(request.Body)
		}

		return req, logAttrs, nil
	}

	newBody, err := re.customRequest.Body.Transform(rawRequestData)
	if err != nil {
		return nil, logAttrs, fmt.Errorf("failed to transform request body: %w", err)
	}

	newBodyBytes, err := json.Marshal(newBody)
	if err != nil {
		return nil, logAttrs, fmt.Errorf("failed to encode transformed body: %w", err)
	}

	req.SetBody(io.NopCloser(bytes.NewReader(newBodyBytes)))

	return req, logAttrs, nil
}

func (re *RESTHandler) transformResponse(
	resp *http.Response,
) (*http.Response, []slog.Attr, error) {
	if re.customResponse == nil || re.customResponse.Body == nil ||
		re.customResponse.Body.IsZero() {
		return resp, nil, nil
	}

	var responseBody any

	if resp.Body != nil {
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		goutils.CatchWarnErrorFunc(resp.Body.Close)

		if err != nil {
			return resp, nil, fmt.Errorf("failed to decode http response: %w", err)
		}
	}

	responseLogAttrs := make([]slog.Attr, 0, 2)
	responseLogAttrs = append(responseLogAttrs, slog.Any("original_body", responseBody))

	transformedBody, err := re.customResponse.Body.Transform(responseBody)
	if err != nil {
		return nil, nil, err
	}

	responseLogAttrs = append(responseLogAttrs, slog.Any("response_body", transformedBody))

	buf := new(bytes.Buffer)

	err = json.NewEncoder(buf).Encode(transformedBody)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to decode transformed response: %w", err)
	}

	resp.Body = io.NopCloser(buf)

	return resp, responseLogAttrs, err
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
		slog.String("handler_type", "rest"),
		slog.String("request_url", request.URL.String()),
		otelutils.NewHeaderLogGroupAttrs(
			"request_headers",
			otelutils.NewTelemetryHeaders(request.Header),
		),
	)

	logger.LogAttrs(ctx, slog.LevelDebug, message, attrs...)
}
