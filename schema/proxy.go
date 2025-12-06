package schema

import (
	"context"
	"net/http"

	"github.com/hasura/goenvconf"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/gohttpc/loadbalancer"
	"github.com/relychan/gotransform"
	wk8orderedmap "github.com/wk8/go-ordered-map/v2"
)

// RelyProxyHandleOptions hold request options for the proxy handler.
type RelyProxyHandleOptions struct {
	HTTPClient     *loadbalancer.LoadBalancerClient
	Settings       *RelyProxySettings
	ParamValues    map[string]string
	DefaultHeaders map[string]string
}

// RelyProxyHandler abstracts the executor to proxy HTTP requests.
type RelyProxyHandler interface {
	// Type returns type of the current handler.
	Type() RelyProxyType
	// Handle resolves the HTTP request and proxies that request to the remote server.
	Handle(
		ctx context.Context,
		request *http.Request,
		options *RelyProxyHandleOptions,
	) (*http.Response, any, error)
}

// NewRelyProxyHandlerOptions hold request options for the proxy handler.
type NewRelyProxyHandlerOptions struct {
	Method     string
	Parameters []Parameter
}

// NewRelyProxyHandlerFunc abstracts a function to create a new proxy handler.
type NewRelyProxyHandlerFunc func(operation *RelyProxyOperation, options *NewRelyProxyHandlerOptions) (RelyProxyHandler, error)

// RelyProxyAction represents a proxy action.
type RelyProxyAction struct {
	// Type of the proxy action.
	Type RelyProxyType `json:"type" yaml:"type"`
	// Overrides the request path. Use the original request path if empty.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Configurations for the proxy request.
	Request RelyProxyGraphQLRequestConfig `json:"request" yaml:"request"`
	// Configurations for evaluating graphql responses.
	Response RelyProxyGraphQLResponseConfig `json:"response" yaml:"response"`
}

// JSONSchema defines a custom definition for JSON schema.
func (RelyProxyAction) JSONSchema() *jsonschema.Schema {
	restSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	restSchema.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Type of the proxy action",
		Enum:        []any{ProxyTypeREST},
	})
	restSchema.Set("path", &jsonschema.Schema{
		Description: "Overrides the request path. Use the original request path if empty",
		Type:        "string",
	})

	graphqlSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	graphqlSchema.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Type of the proxy action",
		Enum:        []any{ProxyTypeGraphQL},
	})
	graphqlSchema.Set("request", &jsonschema.Schema{
		Description: "Configuration for the GraphQL request",
		Ref:         "#/$defs/RelyProxyGraphQLRequestConfig",
	})
	graphqlSchema.Set("response", &jsonschema.Schema{
		Description: "Configuration for the GraphQL response",
		Ref:         "#/$defs/RelyProxyGraphQLResponseConfig",
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "object",
				Description: "Proxy configuration to the remote REST service",
				Required:    []string{"type"},
				Properties:  restSchema,
			},
			{
				Type:        "object",
				Description: "Configurations for proxying request to the remote GraphQL server",
				Properties:  graphqlSchema,
				Required:    []string{"type", "request"},
			},
		},
	}
}

// GraphQLVariableDefinition defines information of the GraphQL variable.
type GraphQLVariableDefinition struct {
	Expression string            `json:"expression,omitempty" yaml:"expression,omitempty"`
	Default    *goenvconf.EnvAny `json:"default,omitempty" yaml:"default,omitempty"`
}

// RelyProxyGraphQLRequestConfig represents configurations for the proxy request.
type RelyProxyGraphQLRequestConfig struct {
	// GraphQL query
	Query string `json:"query,omitempty" yaml:"query,omitempty"`
	// Definition of GraphQL variables.
	Variables *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"variables,omitempty" yaml:"variables,omitempty"`
	// Definition of GraphQL extensions.
	Extensions *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// JSONSchema defines a custom definition for JSON schema.
func (RelyProxyGraphQLRequestConfig) JSONSchema() *jsonschema.Schema {
	graphqlProps := wk8orderedmap.New[string, *jsonschema.Schema]()

	graphqlProps.Set("query", &jsonschema.Schema{
		Description: "GraphQL query string to send",
		Type:        "string",
	})
	graphqlProps.Set("variables", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/GraphQLVariableDefinition",
		},
	})
	graphqlProps.Set("extensions", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/GraphQLVariableDefinition",
		},
	})

	return &jsonschema.Schema{
		Type:       "object",
		Properties: graphqlProps,
		Required:   []string{"query"},
	}
}

// RelyProxyGraphQLResponseConfig represents configurations for the proxy response.
type RelyProxyGraphQLResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int                                   `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"min=400,max=599,default=400"`
	Transform     *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty" jsonschema:"oneof_ref=https://raw.githubusercontent.com/relychan/gotransform/refs/heads/main/jsonschema/gotransform.schema.json,oneof_type=null"` //nolint:lll
}

// IsZero checks if the configuration is empty.
func (conf RelyProxyGraphQLResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil && conf.Transform == nil
}
