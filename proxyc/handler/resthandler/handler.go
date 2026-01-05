// Package resthandler evaluates and execute REST requests to the remote server.
package resthandler

import (
	"context"
	"net/http"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
)

// RESTHandler implements the RelixyHandler interface for REST proxy.
type RESTHandler struct {
	method      string
	requestPath string
	parameters  []*highv3.Parameter
}

// NewRESTHandler creates a RESTHandler from operation.
func NewRESTHandler( //nolint:ireturn
	operation *highv3.Operation,
	proxyAction *base_schema.RelixyAction,
	options *proxyhandler.NewRelixyHandlerOptions,
) (proxyhandler.RelixyHandler, error) {
	handler := &RESTHandler{
		method:     options.Method,
		parameters: openapi.MergeParameters(options.Parameters, operation.Parameters),
	}

	if proxyAction != nil && proxyAction.Path != "" {
		handler.requestPath = proxyAction.Path
	}

	return handler, nil
}

// Type returns type of the current handler.
func (*RESTHandler) Type() base_schema.RelixyType {
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

	req := options.NewRequest(re.method, requestPath)

	if request.Body != nil {
		req.SetBody(request.Body)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return resp, nil, err
	}

	return resp, resp.Body, err
}
