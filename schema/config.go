package schema

import (
	"github.com/hasura/goenvconf"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/gohttpc/loadbalancer"
)

// RelyProxySettings hold settings of the rely proxy.
type RelyProxySettings struct {
	// Set the base path for all API handlers.
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	// Global settings for the HTTP client.
	HTTP *httpconfig.HTTPClientConfig `json:"http,omitempty" yaml:"http,omitempty"`
	// Headers define custom headers to be injected to the remote server.
	// Merged with the global headers.
	Headers map[string]goenvconf.EnvString `json:"headers,omitempty" yaml:"headers,omitempty"`
	// ForwardHeaders define configurations for headers forwarding
	ForwardHeaders RelyProxyForwardHeadersConfig `json:"forwardHeaders,omitempty" yaml:"forwardHeaders,omitempty"`
	// SecuritySchemes define security schemes that can be used by the operations.
	SecuritySchemes []RelyProxySecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// HealthCheck define the health check policy for load balancer recovery.
	HealthCheck *RelyProxyHealthCheckConfig `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
}

// RelyProxyHealthCheckConfig holds health check configurations for server recovery.
type RelyProxyHealthCheckConfig struct {
	// Configurations for health check through HTTP protocol.
	HTTP *loadbalancer.HTTPHealthCheckConfig `json:"http,omitempty" yaml:"http,omitempty"`
}

// RelyProxyForwardHeadersConfig contains configurations for headers forwarding,.
type RelyProxyForwardHeadersConfig struct {
	// Defines header names to be forwarded from the client request.
	Request []string `json:"request,omitempty" yaml:"request,omitempty"`
	// Defines header names to be forwarded from the response.
	Response []string `json:"response,omitempty" yaml:"response,omitempty"`
}

// RelyProxyGraphQLConfig contains configurations for GraphQL proxy.
type RelyProxyGraphQLConfig struct {
	// ScalarTypeMapping configures the custom type mapping between GraphQL scalar types and primitive types for data conversion.
	// Default scalar types are supported. Other types which aren't configured will be forwarded directly
	// from the request parameters and body without serialization.
	ScalarTypeMapping map[string]PrimitiveType `json:"scalarTypeMapping,omitempty" yaml:"scalarTypeMapping,omitempty"`
}
