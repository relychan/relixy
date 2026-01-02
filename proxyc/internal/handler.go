package internal

import (
	"errors"
	"fmt"

	"github.com/relychan/relixy/proxyc/handler/graphqlhandler"
	"github.com/relychan/relixy/proxyc/handler/resthandler"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
)

var (
	ErrUnsupportedProxyType = errors.New("unsupported proxy type")
	ErrWildcardMustBeLast   = errors.New(
		"wildcard '*' must be the last value in a route. trim trailing text or use a '{param}' instead",
	)
	ErrMissingClosingBracket            = errors.New("route param closing delimiter '}' is missing")
	ErrDuplicatedParamKey               = errors.New("routing pattern contains duplicate param key")
	ErrInvalidRegexpPatternParamInRoute = errors.New("invalid regexp pattern in route param")
	ErrReplaceMissingChildNode          = errors.New("replacing missing child node")
)

var proxyHandlerConstructors = map[base_schema.RelixyType]openapi.NewRelixyHandlerFunc{
	base_schema.ProxyTypeREST:    resthandler.NewRESTHandler,
	base_schema.ProxyTypeGraphQL: graphqlhandler.NewGraphQLHandler,
}

// NewProxyHandler creates a proxy handler by type.
func NewProxyHandler( //nolint:ireturn
	operation *openapi.RelixyOpenAPIv3Operation,
	options *openapi.NewRelixyHandlerOptions,
) (openapi.RelixyHandler, error) {
	proxyType := base_schema.ProxyTypeREST

	if operation.Proxy.Type != "" {
		proxyType = operation.Proxy.Type
	}

	constructor, ok := proxyHandlerConstructors[proxyType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProxyType, proxyType)
	}

	return constructor(operation, options)
}

// RegisterProxyHandler registers the handler to the global registry.
func RegisterProxyHandler(
	proxyType base_schema.RelixyType,
	constructor openapi.NewRelixyHandlerFunc,
) {
	proxyHandlerConstructors[proxyType] = constructor
}
