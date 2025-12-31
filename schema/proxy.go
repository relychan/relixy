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

// RelixyHandleOptions hold request options for the proxy handler.
type RelixyHandleOptions struct {
	HTTPClient     *loadbalancer.LoadBalancerClient
	Settings       *RelixySettings
	Path           string
	ParamValues    map[string]string
	DefaultHeaders map[string]string
}

// RelixyHandler abstracts the executor to proxy HTTP requests.
type RelixyHandler interface {
	// Type returns type of the current handler.
	Type() RelixyType
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
type NewRelixyHandlerFunc func(operation *RelixyOperation, options *NewRelixyHandlerOptions) (RelixyHandler, error)

// RelixyAction represents a proxy action.
type RelixyAction struct {
	// Type of the proxy action.
	Type RelixyType `json:"type" yaml:"type"`
	// Overrides the request path. Use the original request path if empty.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Configurations for the proxy request.
	Request RelixyGraphQLRequestConfig `json:"request" yaml:"request"`
	// Configurations for evaluating graphql responses.
	Response RelixyGraphQLResponseConfig `json:"response" yaml:"response"`
}

// JSONSchema defines a custom definition for JSON schema.
func (RelixyAction) JSONSchema() *jsonschema.Schema {
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
		Ref:         "#/$defs/RelixyGraphQLRequestConfig",
	})
	graphqlSchema.Set("response", &jsonschema.Schema{
		Description: "Configuration for the GraphQL response",
		Ref:         "#/$defs/RelixyGraphQLResponseConfig",
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "object",
				Title:       "RelixyActionREST",
				Description: "Proxy configuration to the remote REST service",
				Required:    []string{"type"},
				Properties:  restSchema,
			},
			{
				Type:        "object",
				Title:       "RelixyActionGraphQL",
				Description: "Configurations for proxying request to the remote GraphQL server",
				Properties:  graphqlSchema,
				Required:    []string{"type", "request"},
			},
		},
	}
}

// GraphQLVariableDefinition defines information of the GraphQL variable.
type GraphQLVariableDefinition struct {
	// JMESPath to evaluate the variable from request.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Default value if the path or value is empty.
	Default *goenvconf.EnvAny `json:"default,omitempty" yaml:"default,omitempty"`
}

// RelixyGraphQLRequestConfig represents configurations for the proxy request.
type RelixyGraphQLRequestConfig struct {
	// GraphQL query
	Query string `json:"query,omitempty" yaml:"query,omitempty"`
	// Definition of GraphQL variables.
	Variables *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"variables,omitempty" yaml:"variables,omitempty"`
	// Definition of GraphQL extensions.
	Extensions *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// JSONSchema defines a custom definition for JSON schema.
func (RelixyGraphQLRequestConfig) JSONSchema() *jsonschema.Schema {
	graphqlProps := wk8orderedmap.New[string, *jsonschema.Schema]()

	graphqlProps.Set("query", &jsonschema.Schema{
		Description: "GraphQL query string to send",
		Type:        "string",
	})
	graphqlProps.Set("variables", &jsonschema.Schema{
		Type:        "object",
		Description: "Definition of GraphQL variables",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/GraphQLVariableDefinition",
		},
	})
	graphqlProps.Set("extensions", &jsonschema.Schema{
		Type:        "object",
		Description: "Definition of GraphQL extensions",
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

// RelixyGraphQLResponseConfig represents configurations for the proxy response.
type RelixyGraphQLResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"min=400,max=599,default=400"`
	// Configurations for transforming response data.
	Transform *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty" jsonschema:"oneof_ref=https://raw.githubusercontent.com/relychan/gotransform/refs/heads/main/jsonschema/gotransform.schema.json,oneof_type=null"` //nolint:lll
}

// IsZero checks if the configuration is empty.
func (conf RelixyGraphQLResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Transform == nil || conf.Transform.IsZero())
}
