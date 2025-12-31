// Package resthandler evaluates and execute REST requests to the remote server.
package resthandler

import (
	"context"
	"net/http"

	"github.com/relychan/relixy/schema"
)

// RESTHandler implements the RelixyHandler interface for REST proxy.
type RESTHandler struct {
	method      string
	requestPath string
	parameters  []schema.Parameter
}

// NewRESTHandler creates a RESTHandler from operation.
func NewRESTHandler( //nolint:ireturn
	operation *schema.RelixyOperation,
	options *schema.NewRelixyHandlerOptions,
) (schema.RelixyHandler, error) {
	return &RESTHandler{
		method:      options.Method,
		requestPath: operation.Proxy.Path,
		parameters:  schema.MergeParameters(options.Parameters, operation.Parameters),
	}, nil
}

// Type returns type of the current handler.
func (*RESTHandler) Type() schema.RelixyType {
	return schema.ProxyTypeREST
}

// Handle resolves the HTTP request and proxies that request to the remote server.
func (re *RESTHandler) Handle(
	ctx context.Context,
	request *http.Request,
	options *schema.RelixyHandleOptions,
) (*http.Response, any, error) {
	requestPath := re.requestPath
	if requestPath == "" {
		requestPath = options.Path
	}

	if request.URL.RawQuery != "" {
		requestPath += "?" + request.URL.RawQuery
	}

	req := options.HTTPClient.R(re.method, requestPath)

	if options.Settings.ForwardHeaders != nil {
		for _, key := range options.Settings.ForwardHeaders.Request {
			value := req.Header().Get(key)
			if value != "" {
				req.Header().Set(key, value)
			}
		}
	}

	if request.Body != nil {
		req.SetBody(request.Body)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return resp, nil, err
	}

	return resp, resp.Body, err
}
