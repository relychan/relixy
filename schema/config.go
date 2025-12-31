package schema

import (
	"github.com/hasura/goenvconf"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/gohttpc/loadbalancer"
)

// RelixySettings hold settings of the rely proxy.
type RelixySettings struct {
	// Global settings for the HTTP client.
	HTTP *httpconfig.HTTPClientConfig `json:"http,omitempty" yaml:"http,omitempty"`
	// Headers define custom headers to be injected to the remote server.
	// Merged with the global headers.
	Headers map[string]goenvconf.EnvString `json:"headers,omitempty" yaml:"headers,omitempty"`
	// ForwardHeaders define configurations for headers forwarding
	ForwardHeaders *RelixyForwardHeadersConfig `json:"forwardHeaders,omitempty" yaml:"forwardHeaders,omitempty"`
	// SecuritySchemes define security schemes that can be used by the operations.
	SecuritySchemes []RelixySecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// HealthCheck define the health check policy for load balancer recovery.
	HealthCheck *RelixyHealthCheckConfig `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
}

// RelixyHealthCheckConfig holds health check configurations for server recovery.
type RelixyHealthCheckConfig struct {
	// Configurations for health check through HTTP protocol.
	HTTP *loadbalancer.HTTPHealthCheckConfig `json:"http,omitempty" yaml:"http,omitempty"`
}

// RelixyForwardHeadersConfig contains configurations for headers forwarding,.
type RelixyForwardHeadersConfig struct {
	// Defines header names to be forwarded from the client request.
	Request []string `json:"request,omitempty" yaml:"request,omitempty"`
	// Defines header names to be forwarded from the response.
	Response []string `json:"response,omitempty" yaml:"response,omitempty"`
}

// RelixyGraphQLConfig contains configurations for GraphQL proxy.
type RelixyGraphQLConfig struct {
	// ScalarTypeMapping configures the custom type mapping between GraphQL scalar types and primitive types for data conversion.
	// Default scalar types are supported. Other types which aren't configured will be forwarded directly
	// from the request parameters and body without serialization.
	ScalarTypeMapping map[string]PrimitiveType `json:"scalarTypeMapping,omitempty" yaml:"scalarTypeMapping,omitempty"`
}
