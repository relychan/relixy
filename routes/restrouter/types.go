package restrouter

import (
	"github.com/relychan/relixy/proxyc"
)

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
