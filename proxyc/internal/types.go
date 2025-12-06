package internal

import (
	"errors"
	"fmt"

	"github.com/relychan/relyx/proxyc/handler/graphqlhandler"
	"github.com/relychan/relyx/proxyc/handler/resthandler"
	"github.com/relychan/relyx/schema"
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

var proxyHandlerConstructors = map[schema.RelyProxyType]schema.NewRelyProxyHandlerFunc{
	schema.ProxyTypeREST:    resthandler.NewRESTHandler,
	schema.ProxyTypeGraphQL: graphqlhandler.NewGraphQLHandler,
}

// NewProxyHandler creates a proxy handler by type.
func NewProxyHandler( //nolint:ireturn
	operation *schema.RelyProxyOperation,
	options *schema.NewRelyProxyHandlerOptions,
) (schema.RelyProxyHandler, error) {
	proxyType := schema.ProxyTypeREST

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
	proxyType schema.RelyProxyType,
	constructor schema.NewRelyProxyHandlerFunc,
) {
	proxyHandlerConstructors[proxyType] = constructor
}
