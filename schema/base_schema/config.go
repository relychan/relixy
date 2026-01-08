package base_schema

import (
	"github.com/relychan/gohttpc/loadbalancer"
	"github.com/relychan/gotransform"
)

// RelixyHealthCheckConfig holds health check configurations for server recovery.
type RelixyHealthCheckConfig struct {
	// Configurations for health check through HTTP protocol.
	HTTP *loadbalancer.HTTPHealthCheckConfig `json:"http,omitempty" yaml:"http,omitempty"`
}

// RelixyResponseConfig represents configurations for the proxy response.
type RelixyResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"minimum=400,maximum=599,default=400"`
	// Configurations for transforming response data.
	Transform *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty"`
}

// IsZero checks if the configuration is empty.
func (conf RelixyResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Transform == nil || conf.Transform.IsZero())
}
