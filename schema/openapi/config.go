// Package openapi defines schemas for Open API resources
package openapi

import (
	"github.com/hasura/goenvconf"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/relixy/schema/base_schema"
)

// RelixyOpenAPISettings hold settings of the rely proxy.
type RelixyOpenAPISettings struct {
	// Base path of the resource.
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	// Global settings for the HTTP client.
	HTTP *httpconfig.HTTPClientConfig `json:"http,omitempty" yaml:"http,omitempty"`
	// Headers define custom headers to be injected to the remote server.
	// Merged with the global headers.
	Headers map[string]goenvconf.EnvString `json:"headers,omitempty" yaml:"headers,omitempty"`
	// ForwardHeaders define configurations for headers forwarding
	ForwardHeaders *RelixyOpenAPIForwardHeadersConfig `json:"forwardHeaders,omitempty" yaml:"forwardHeaders,omitempty"`
	// HealthCheck define the health check policy for load balancer recovery.
	HealthCheck *base_schema.RelixyHealthCheckConfig `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
}

// RelixyOpenAPIForwardHeadersConfig contains configurations for headers forwarding,.
type RelixyOpenAPIForwardHeadersConfig struct {
	// Defines header names to be forwarded from the client request.
	Request []string `json:"request,omitempty" yaml:"request,omitempty"`
	// Defines header names to be forwarded from the response.
	Response []string `json:"response,omitempty" yaml:"response,omitempty"`
}
