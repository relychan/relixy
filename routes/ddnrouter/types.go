package ddnrouter

import (
	"errors"
	"strings"

	"github.com/relychan/relixy/proxyc"
)

// State holds common states of the handler.
type State struct {
	ProxyClients []*proxyc.ProxyClient
}

// Close stops internal processes of the state.
func (s *State) FindProxyClient(requestPath string) *proxyc.ProxyClient {
	defaultIndex := -1

	for i, pc := range s.ProxyClients {
		metadata := pc.Metadata()

		if metadata.Definition.Settings == nil ||
			metadata.Definition.Settings.BasePath == "" ||
			metadata.Definition.Settings.BasePath == "/" {
			defaultIndex = i

			continue
		}

		if strings.HasPrefix(requestPath, metadata.Definition.Settings.BasePath) {
			return pc
		}
	}

	if defaultIndex > -1 {
		return s.ProxyClients[defaultIndex]
	}

	return nil
}

// Close stops internal processes of the state.
func (s *State) Close() error {
	var errs []error

	for _, pc := range s.ProxyClients {
		err := pc.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
