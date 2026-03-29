package restrouter

import (
	"github.com/relychan/openapitools/openapiclient"
)

// State holds common states of the handler.
type State struct {
	ProxyClient *openapiclient.ProxyClient
}

// Close stops internal processes of the state.
func (s *State) Close() error {
	if s.ProxyClient != nil {
		return s.ProxyClient.Close()
	}

	return nil
}
