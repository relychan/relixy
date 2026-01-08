// Package handler defines the global proxy handler with default constructors
package handler

import (
	"errors"
	"fmt"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/graphqlhandler"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/proxyc/handler/resthandler"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
)

var ErrUnsupportedProxyType = errors.New("unsupported proxy type")

var proxyHandlerConstructors = map[base_schema.RelixyActionType]proxyhandler.NewRelixyHandlerFunc{
	base_schema.ProxyTypeREST:    resthandler.NewRESTHandler,
	base_schema.ProxyTypeGraphQL: graphqlhandler.NewGraphQLHandler,
}

// NewProxyHandler creates a proxy handler by type.
func NewProxyHandler( //nolint:ireturn,nolintlint
	operation *highv3.Operation,
	options *proxyhandler.NewRelixyHandlerOptions,
) (proxyhandler.RelixyHandler, error) {
	var proxyAction base_schema.RelixyAction

	if operation.Extensions != nil {
		rawProxyAction, exist := operation.Extensions.Get(openapi.XRelyProxyAction)
		if exist && rawProxyAction != nil {
			err := rawProxyAction.Decode(&proxyAction)
			if err != nil {
				return nil, err
			}
		}
	}

	switch proxyAction.Type {
	case base_schema.ProxyTypeGraphQL:
		_, err := graphqlhandler.ValidateGraphQLString(proxyAction.Request.Query)
		if err != nil {
			return nil, err
		}
	case "":
		proxyAction.Type = base_schema.ProxyTypeREST
	default:
	}

	constructor, ok := proxyHandlerConstructors[proxyAction.Type]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProxyType, proxyAction.Type)
	}

	return constructor(operation, &proxyAction, options)
}

// RegisterProxyHandler registers the handler to the global registry.
func RegisterProxyHandler(
	proxyType base_schema.RelixyActionType,
	constructor proxyhandler.NewRelixyHandlerFunc,
) {
	proxyHandlerConstructors[proxyType] = constructor
}
