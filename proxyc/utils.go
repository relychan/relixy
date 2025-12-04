package proxyc

import (
	"fmt"
	"net/url"
	"slices"
)

// ParseHTTPURL parses and validates the HTTP URL.
func ParseHTTPURL(input string, bannedHosts []string) (*url.URL, error) {
	result, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	if result.Scheme != "http" && result.Scheme != "https" {
		return nil, fmt.Errorf("%w: %s", errInvalidHTTPURL, input)
	}

	// ban localhost and suspicious hosts
	if len(bannedHosts) > 0 && slices.Contains(bannedHosts, result.Hostname()) {
		return nil, fmt.Errorf("%w: %s", errForbiddenHost, result.Host)
	}

	return result, nil
}
