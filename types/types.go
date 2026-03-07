// Package types defines common types for the rely proxy service.
package types //nolint:revise

import (
	"github.com/relychan/relixy/proxyc"
)

// BuildVersion is set when building the binary.
var BuildVersion = "dev"

// State holds common states of the handler.
type State struct {
	ProxyClient *proxyc.ProxyClient
}

// Close stops internal processes of the state.
func (s *State) Close() error {
	if s.ProxyClient != nil {
		return s.ProxyClient.Close()
	}

	return nil
}
