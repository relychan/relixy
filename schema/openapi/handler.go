package openapi

import (
	"context"
	"net/http"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gohttpc/loadbalancer"
	"github.com/relychan/relixy/schema/base_schema"
)

// RelixyHandleOptions hold request options for the proxy handler.
type RelixyHandleOptions struct {
	HTTPClient     *loadbalancer.LoadBalancerClient
	Settings       *RelixyOpenAPISettings
	Path           string
	ParamValues    map[string]string
	DefaultHeaders map[string]string
}

// RelixyHandler abstracts the executor to proxy HTTP requests.
type RelixyHandler interface {
	// Type returns type of the current handler.
	Type() base_schema.RelixyType
	// Handle resolves the HTTP request and proxies that request to the remote server.
	Handle(
		ctx context.Context,
		request *http.Request,
		options *RelixyHandleOptions,
	) (*http.Response, any, error)
}

// NewRelixyHandlerOptions hold request options for the proxy handler.
type NewRelixyHandlerOptions struct {
	Method     string
	Parameters []Parameter
	GetEnv     goenvconf.GetEnvFunc
}

// GetEnvFunc returns a function to get environment variables.
func (nrp NewRelixyHandlerOptions) GetEnvFunc() goenvconf.GetEnvFunc {
	if nrp.GetEnv == nil {
		return goenvconf.GetOSEnv
	}

	return nrp.GetEnv
}

// NewRelixyHandlerFunc abstracts a function to create a new proxy handler.
type NewRelixyHandlerFunc func(operation *RelixyOpenAPIv3Operation, options *NewRelixyHandlerOptions) (RelixyHandler, error)
