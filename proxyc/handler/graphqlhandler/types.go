package graphqlhandler

import (
	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
	"github.com/relychan/gotransform/jmes"
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

// RelixyGraphQLActionConfig represents a proxy action config for GraphQL.
type RelixyGraphQLActionConfig struct {
	// Type of the proxy action which is always graphql
	Type proxyhandler.ProxyActionType `json:"type" yaml:"type" jsonschema:"enum=graphql"`
	// Configurations for the GraphQL proxy request.
	Request *RelixyGraphQLRequestConfig `json:"request,omitempty" yaml:"request,omitempty"`
	// Configurations for evaluating graphql responses.
	Response *RelixyCustomGraphQLResponseConfig `json:"response,omitempty" yaml:"response,omitempty"`
}

// RelixyGraphQLRequestConfig represents configurations for the proxy request.
type RelixyGraphQLRequestConfig struct {
	// GraphQL query
	Query string `json:"query" yaml:"query"`
	// Definition of GraphQL variables.
	Variables map[string]jmes.FieldMappingEntryConfig `json:"variables,omitempty" yaml:"variables,omitempty"`
	// Definition of GraphQL extensions.
	Extensions map[string]jmes.FieldMappingEntryConfig `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// RelixyCustomGraphQLResponseConfig represents configurations for the proxy response.
type RelixyCustomGraphQLResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the upstream response which is usually 200 OK.
	HTTPErrorCode *int `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"minimum=400,maximum=599,default=400"`
	// Configurations for transforming response data.
	Body *gotransform.TemplateTransformerConfig `json:"body,omitempty" yaml:"body,omitempty"`
}

// IsZero checks if the configuration is empty.
func (conf RelixyCustomGraphQLResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Body == nil || conf.Body.IsZero())
}

// RelixyCustomGraphQLResponse represents configurations for the proxy response.
type RelixyCustomGraphQLResponse struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the upstream response which is usually 200 OK.
	HTTPErrorCode *int
	// Configurations for transforming response body data.
	Body gotransform.TemplateTransformer
}

// NewRelixyCustomGraphQLResponse creates a [RelixyCustomGraphQLResponse] from raw configurations.
func NewRelixyCustomGraphQLResponse(
	config *RelixyCustomGraphQLResponseConfig,
	getEnv goenvconf.GetEnvFunc,
) (*RelixyCustomGraphQLResponse, error) {
	if config == nil || config.IsZero() {
		return nil, nil
	}

	result := &RelixyCustomGraphQLResponse{
		HTTPErrorCode: config.HTTPErrorCode,
	}

	if config.Body != nil {
		transformer, err := gotransform.NewTransformerFromConfig("", *config.Body, getEnv)
		if err != nil {
			return result, err
		}

		result.Body = transformer
	}

	return result, nil
}

// IsZero checks if the configuration is empty.
func (conf RelixyCustomGraphQLResponse) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Body == nil || conf.Body.IsZero())
}
