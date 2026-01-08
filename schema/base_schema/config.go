package base_schema

import (
	"github.com/relychan/gohttpc/loadbalancer"
)

// RelixyHealthCheckConfig holds health check configurations for server recovery.
type RelixyHealthCheckConfig struct {
	// Configurations for health check through HTTP protocol.
	HTTP *loadbalancer.HTTPHealthCheckConfig `json:"http,omitempty" yaml:"http,omitempty"`
}
