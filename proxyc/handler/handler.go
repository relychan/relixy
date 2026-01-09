// Package handler defines the global proxy handler with default constructors
package handler

import (
	"errors"
	"fmt"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/graphqlhandler"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/proxyc/handler/resthandler"
	"github.com/relychan/relixy/schema/openapi"
	"go.yaml.in/yaml/v4"
)

var ErrUnsupportedProxyType = errors.New("unsupported proxy type")

var proxyHandlerConstructors = map[proxyhandler.ProxyActionType]proxyhandler.NewRelixyHandlerFunc{
	resthandler.ProxyActionTypeREST: resthandler.NewRESTHandler,
	graphqlhandler.ProxyTypeGraphQL: graphqlhandler.NewGraphQLHandler,
}

// NewProxyHandler creates a proxy handler by type.
func NewProxyHandler( //nolint:ireturn,nolintlint
	operation *highv3.Operation,
	options *proxyhandler.NewRelixyHandlerOptions,
) (proxyhandler.RelixyHandler, error) {
	var proxyAction rawRelixyActionConfig

	var rawProxyAction *yaml.Node

	if operation.Extensions != nil {
		var exist bool

		rawProxyAction, exist = operation.Extensions.Get(openapi.XRelyProxyAction)
		if exist && rawProxyAction != nil {
			err := rawProxyAction.Decode(&proxyAction)
			if err != nil {
				return nil, err
			}
		}
	}

	if proxyAction.Type == "" {
		proxyAction.Type = resthandler.ProxyActionTypeREST
	}

	constructor, ok := proxyHandlerConstructors[proxyAction.Type]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProxyType, proxyAction.Type)
	}

	return constructor(operation, rawProxyAction, options)
}

// RegisterProxyHandler registers the handler to the global registry.
func RegisterProxyHandler(
	proxyType proxyhandler.ProxyActionType,
	constructor proxyhandler.NewRelixyHandlerFunc,
) {
	proxyHandlerConstructors[proxyType] = constructor
}

// rawRelixyActionConfig represents a raw proxy action with type only.
type rawRelixyActionConfig struct {
	// Type of the proxy action.
	Type proxyhandler.ProxyActionType `json:"type" yaml:"type"`
}
