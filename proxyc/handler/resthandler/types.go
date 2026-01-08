package resthandler

import (
	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
)

// ProxyActionTypeREST represents a constant value for REST proxy action.
const ProxyActionTypeREST proxyhandler.ProxyActionType = "rest"

// RelixyRESTActionConfig represents a proxy action config for REST operation.
type RelixyRESTActionConfig struct {
	// Type of the proxy action which is always graphql
	Type proxyhandler.ProxyActionType `json:"type" yaml:"type" jsonschema:"enum=rest"`
	// Configurations for the GraphQL proxy request.
	Request *RelixyRESTRequestConfig `json:"request" yaml:"request"`
	// Configurations for evaluating graphql responses.
	Response *proxyhandler.RelixyResponseRawConfig `json:"response" yaml:"response"`
}

// GraphQLVariableDefinition defines information of the GraphQL variable.
type GraphQLVariableDefinition struct {
	// JMESPath to evaluate the variable from request.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Default value if the path or value is empty.
	Default *goenvconf.EnvAny `json:"default,omitempty" yaml:"default,omitempty"`
}

// RelixyRESTRequestConfig represents configurations for the proxy request.
type RelixyRESTRequestConfig struct {
	// Overrides the request path. Use the original request path if empty.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Definition of request body template.
	Body *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty"`
}
