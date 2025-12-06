// Package graphqlhandler evaluates and execute GraphQL requests to the remote server.
package graphqlhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/jmespath-community/go-jmespath"
	"github.com/relychan/goutils"
	"github.com/relychan/goutils/httpheader"
	"github.com/relychan/relyx/schema"
	"github.com/vektah/gqlparser/ast"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GraphQLHandler implements the RelyProxyHandler interface for GraphQL proxy.
type GraphQLHandler struct {
	requestPath         string
	parameters          []schema.Parameter
	query               string
	operation           ast.Operation
	variableDefinitions ast.VariableDefinitionList
	variables           map[string]graphqlVariable
	extensions          map[string]graphqlVariable
	operationName       string
	responseConfig      schema.RelyProxyGraphQLResponseConfig
}

// NewGraphQLHandler creates a GraphQL request from operation.
func NewGraphQLHandler( //nolint:ireturn,nolintlint
	operation *schema.RelyProxyOperation,
	options *schema.NewRelyProxyHandlerOptions,
) (schema.RelyProxyHandler, error) {
	handler, err := ValidateGraphQLString(operation.Proxy.Request.Query)
	if err != nil {
		return nil, err
	}

	handler.requestPath = operation.Proxy.Path
	handler.responseConfig = operation.Proxy.Response
	handler.parameters = schema.MergeParameters(options.Parameters, operation.Parameters)

	handler.variables, err = validateGraphQLVariables(operation.Proxy.Request.Variables)
	if err != nil {
		return nil, err
	}

	handler.extensions, err = validateGraphQLVariables(operation.Proxy.Request.Extensions)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// Type returns type of the current handler.
func (*GraphQLHandler) Type() schema.RelyProxyType {
	return schema.ProxyTypeGraphQL
}

// Handle resolves the HTTP request and proxies that request to the remote server.
func (ge *GraphQLHandler) Handle( //nolint:funlen
	ctx context.Context,
	request *http.Request,
	options *schema.RelyProxyHandleOptions,
) (*http.Response, any, error) {
	span := trace.SpanFromContext(ctx)
	requestPath := options.Settings.BasePath

	if ge.requestPath != "" {
		requestPath = ge.requestPath
	}

	graphqlPayload := GraphQLRequestBody{
		Query:         ge.query,
		OperationName: ge.operationName,
	}

	span.SetAttributes(
		attribute.String("graphql.operation.name", ge.operationName),
		attribute.String("graphql.operation.type", string(ge.operation)),
		attribute.String("graphql.query", ge.query),
	)

	logAttrs := []slog.Attr{
		slog.String("path", requestPath),
	}

	requestHeaders := map[string]string{}

	for key, header := range request.Header {
		if len(header) == 0 {
			continue
		}

		requestHeaders[strings.ToLower(key)] = header[0]
	}

	requestData := &requestTemplateData{
		Params:      options.ParamValues,
		QueryParams: request.URL.Query(),
		Headers:     requestHeaders,
	}

	if request.Body != nil {
		var body any

		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			ge.printLog(
				ctx,
				request, "failed to decode body",
				append(
					logAttrs,
					slog.String("error", err.Error()),
				),
			)

			return nil, nil, err
		}

		requestData.Body = body
	}

	variables, err := ge.resolveRequestVariables(requestData)
	if err != nil {
		ge.printLog(
			ctx,
			request, "failed to evaluate variables",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return nil, nil, err
	}

	graphqlPayload.Variables = variables

	graphqlPayload.Extensions, err = ge.resolveRequestExtensions(requestData)
	if err != nil {
		ge.printLog(
			ctx,
			request, "failed to evaluate extensions",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return nil, nil, err
	}

	logAttrs = append(logAttrs,
		slog.GroupAttrs(
			"request_graphql",
			slog.String("query", graphqlPayload.Query),
			slog.Any("variables", graphqlPayload.Variables),
			slog.Any("extensions", graphqlPayload.Extensions),
		),
	)

	req := options.HTTPClient.R(http.MethodPost, requestPath)
	reqHeader := req.Header()

	for key, value := range options.DefaultHeaders {
		reqHeader.Set(key, value)
	}

	if options.Settings.ForwardHeaders != nil {
		for _, key := range options.Settings.ForwardHeaders.Request {
			value := reqHeader.Get(key)
			if value != "" {
				reqHeader.Set(key, value)
			}
		}
	}

	reqHeader.Set(httpheader.ContentType, httpheader.ContentTypeJSON)

	reader := new(bytes.Buffer)

	enc := json.NewEncoder(reader)
	enc.SetEscapeHTML(false)

	err = enc.Encode(graphqlPayload)
	if err != nil {
		ge.printLog(
			ctx, request,
			"failed to encode graphql payload",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return nil, nil, err
	}

	req.SetBody(reader)

	resp, err := req.Execute(ctx)
	if err != nil {
		ge.printLog(
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

	span.SetAttributes(attribute.Int("http.response.original_status_code", resp.StatusCode))

	if ge.responseConfig.IsZero() {
		ge.printLog(
			ctx, request,
			resp.Status,
			logAttrs,
		)

		return resp, resp.Body, nil
	}

	newResp, respBody, respLogAttrs, err := ge.transformResponse(resp)
	logAttrs = append(logAttrs, respLogAttrs...)

	if err != nil {
		ge.printLog(
			ctx, request,
			"failed to transform request",
			append(
				logAttrs,
				slog.String("error", err.Error()),
			),
		)

		return resp, nil, err
	}

	ge.printLog(
		ctx, request,
		resp.Status,
		append(logAttrs, slog.Any("response_body", respBody)),
	)

	return newResp, respBody, err
}

func (ge *GraphQLHandler) transformResponse( //nolint:revive
	resp *http.Response,
) (*http.Response, any, []slog.Attr, error) {
	defer goutils.CatchWarnErrorFunc(resp.Body.Close)

	responseLogAttrs := []slog.Attr{}

	var responseBody map[string]any

	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return resp, nil, responseLogAttrs, fmt.Errorf("failed to decode graphql response: %w", err)
	}

	responseLogAttrs = append(responseLogAttrs, slog.Any("original_body", responseBody))

	if ge.responseConfig.HTTPErrorCode != nil {
		errorBody, hasError := responseBody["errors"]
		if hasError && errorBody != nil {
			// overwrite the error code.
			resp.StatusCode = *ge.responseConfig.HTTPErrorCode
		}
	}

	responseLogAttrs = append(
		responseLogAttrs,
		slog.Int("status_code_final", resp.StatusCode),
	)

	return resp, responseBody, responseLogAttrs, nil
}

func (ge *GraphQLHandler) resolveRequestVariables(
	requestData *requestTemplateData,
) (map[string]any, error) {
	results := make(map[string]any)

	if len(ge.variableDefinitions) == 0 {
		return results, nil
	}

	rawRequestData := requestData.ToMap()

	for _, varDef := range ge.variableDefinitions {
		// Resolve graphql variables. Variables are resolved in order:
		// - In proxy config.
		// - In request parameters, query and body.
		// - Default value in config.
		variable, ok := ge.variables[varDef.Variable]
		if ok {
			if variable.Default != nil {
				results[varDef.Variable] = variable.Default
			}

			if variable.Expression != "" {
				value, err := jmespath.Search(variable.Expression, rawRequestData)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to select value of variable %s: %w",
						varDef.Variable,
						err,
					)
				}

				if value != nil {
					typedValue, err := convertVariableTypeFromUnknownValue(varDef, value)
					if err != nil {
						return nil, fmt.Errorf(
							"failed to evaluate value of variable %s: %w",
							varDef.Variable,
							err,
						)
					}

					results[varDef.Variable] = typedValue
				}

				continue
			}
		}

		if varDef.Variable == "body" {
			results[varDef.Variable] = requestData.Body

			continue
		}

		param, ok := requestData.Params[varDef.Variable]
		if ok && param != "" {
			typedParam, err := convertVariableTypeFromString(varDef, param)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to evaluate the type of variable %s: %w",
					varDef.Variable,
					err,
				)
			}

			results[varDef.Variable] = typedParam

			continue
		}

		queryValue := requestData.QueryParams.Get(varDef.Variable)
		if queryValue != "" {
			typedValue, err := convertVariableTypeFromString(varDef, queryValue)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to evaluate the type of variable %s: %w",
					varDef.Variable,
					err,
				)
			}

			results[varDef.Variable] = typedValue
		}
	}

	return results, nil
}

func (ge *GraphQLHandler) resolveRequestExtensions(
	requestData *requestTemplateData,
) (map[string]any, error) {
	results := make(map[string]any)
	rawRequestData := requestData.ToMap()

	for key, extension := range ge.extensions {
		if extension.Default != nil {
			results[key] = extension.Default
		}

		if extension.Expression != "" {
			value, err := jmespath.Search(extension.Expression, rawRequestData)
			if err != nil {
				return nil, fmt.Errorf("failed to select value of extension %s: %w", key, err)
			}

			if value != nil {
				results[key] = value
			}

			continue
		}
	}

	return results, nil
}

func (ge *GraphQLHandler) printLog(
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
		slog.String("handler_type", string(schema.ProxyTypeGraphQL)),
		slog.String("operation_name", ge.operationName),
		slog.String("operation_type", string(ge.operation)),
		slog.String("request_url", request.URL.String()),
		slog.Any("variable_definitions", ge.variableDefinitions),
		slog.Any("variables", ge.variables),
		otelutils.NewHeaderLogGroupAttrs(
			"request_headers",
			otelutils.NewTelemetryHeaders(request.Header),
		),
	)

	logger.LogAttrs(ctx, slog.LevelDebug, message, attrs...)
}

func convertVariableTypeFromString(varDef *ast.VariableDefinition, value string) (any, error) {
	if varDef.Type == nil {
		// unknown type. Returns the original value.
		return value, nil
	}

	switch strings.ToLower(varDef.Type.NamedType) {
	case "bool", "boolean":
		return strconv.ParseBool(value)
	case "int", "int8", "int16", "int32", "int64":
		return strconv.ParseInt(value, 10, 64)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return strconv.ParseUint(value, 10, 64)
	case "number", "decimal", "float", "float32", "float64", "double":
		return strconv.ParseFloat(value, 64)
	default:
		// unknown type. Returns the original value.
		return value, nil
	}
}

func convertVariableTypeFromUnknownValue(varDef *ast.VariableDefinition, value any) (any, error) {
	if varDef.Type == nil || value == nil {
		// unknown type. Returns the original value.
		return value, nil
	}

	if str, ok := value.(string); ok {
		return convertVariableTypeFromString(varDef, str)
	}

	if strPtr, ok := value.(*string); ok {
		if strPtr == nil {
			return nil, nil
		}

		return convertVariableTypeFromString(varDef, *strPtr)
	}

	switch strings.ToLower(varDef.Type.NamedType) {
	case "bool", "boolean":
		return goutils.DecodeNullableBoolean(value)
	case "int", "int8", "int16", "int32", "int64":
		return goutils.DecodeNullableNumber[int64](value)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return goutils.DecodeNullableNumber[uint64](value)
	case "number", "decimal", "float", "float32", "float64", "double":
		return goutils.DecodeNullableNumber[float64](value)
	default:
		// unknown type. Returns the original value.
		return value, nil
	}
}
