// Package proxyhandler defines types for the proxy handler.
package proxyhandler

import (
	"context"
	"net/http"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/gohttpc"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
)

// RelixyHandleOptions hold request options for the proxy handler.
type RelixyHandleOptions struct {
	NewRequest  NewRequestFunc
	Settings    *openapi.RelixyOpenAPISettings
	Path        string
	ParamValues map[string]string
}

// RelixyHandler abstracts the executor to proxy HTTP requests.
type RelixyHandler interface {
	// Type returns type of the current handler.
	Type() base_schema.RelixyActionType
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
	Parameters []*highv3.Parameter
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
type NewRelixyHandlerFunc func(operation *highv3.Operation, proxyAction *base_schema.RelixyAction, options *NewRelixyHandlerOptions) (RelixyHandler, error)

// NewRequestFunc abstracts a function to create an HTTP request.
type NewRequestFunc func(method string, url string) *gohttpc.RequestWithClient
