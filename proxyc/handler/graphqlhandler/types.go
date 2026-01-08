package graphqlhandler

import (
	"net/url"

	"github.com/hasura/goenvconf"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
)

// ProxyTypeGraphQL represents a constant value for GraphQL proxy action.
const ProxyTypeGraphQL proxyhandler.ProxyActionType = "graphql"

// GraphQLRequestBody represents a request body to a GraphQL server.
type GraphQLRequestBody struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
	Extensions    map[string]any `json:"extensions,omitempty"`
}

type graphqlVariable struct {
	Path    string
	Default any
}

type requestTemplateData struct {
	Params      map[string]string
	QueryParams url.Values
	Headers     map[string]string
	Body        any
}

// ToMap converts the struct to map.
func (rtd requestTemplateData) ToMap() map[string]any {
	result := map[string]any{
		"param":   rtd.Params,
		"query":   rtd.QueryParams,
		"headers": rtd.Headers,
	}

	if rtd.Body != nil {
		result["body"] = rtd.Body
	}

	return result
}

// RelixyGraphQLActionConfig represents a proxy action config for GraphQL.
type RelixyGraphQLActionConfig struct {
	// Type of the proxy action which is always graphql
	Type proxyhandler.ProxyActionType `json:"type" yaml:"type" jsonschema:"enum=graphql"`
	// Configurations for the GraphQL proxy request.
	Request *RelixyGraphQLRequestConfig `json:"request" yaml:"request"`
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

// RelixyGraphQLRequestConfig represents configurations for the proxy request.
type RelixyGraphQLRequestConfig struct {
	// GraphQL query
	Query string `json:"query,omitempty" yaml:"query,omitempty"`
	// Definition of GraphQL variables.
	Variables *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"variables,omitempty" yaml:"variables,omitempty"`
	// Definition of GraphQL extensions.
	Extensions *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyGraphQLRequestConfig) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("variables", &jsonschema.Schema{
			Description: "Definition of GraphQL variables.",
			Type:        "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/GraphQLVariableDefinition",
			},
		})
	schema.Properties.
		Set("extensions", &jsonschema.Schema{
			Description: "Definition of GraphQL extensions.",
			Type:        "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/GraphQLVariableDefinition",
			},
		})
}
