package ddnrouter

import (
	"errors"
	"strings"

	"github.com/relychan/relixy/schema"
)

// State holds common states of the handler.
type State struct {
	ProxyClients []*schema.OpenAPIClient
}

// FindProxyClient find the suitable proxy client from the request path.
func (s *State) FindProxyClient(requestPath string) *schema.OpenAPIClient {
	defaultIndex := -1

	for i, pc := range s.ProxyClients {
		metadata := pc.Metadata()

		if metadata.Settings == nil ||
			metadata.Settings.BasePath == "" ||
			metadata.Settings.BasePath == "/" {
			defaultIndex = i

			continue
		}

		if strings.HasPrefix(requestPath, metadata.Settings.BasePath) {
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
