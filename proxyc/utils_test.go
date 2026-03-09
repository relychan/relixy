package proxyc

import (
	"testing"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"gotest.tools/v3/assert"
)

func TestParseServerURL(t *testing.T) {
	serverURL := "http://localhost:8080"
	t.Setenv("GRAPHQL_SERVER_URL", serverURL)

	testCases := []struct {
		Server   *highv3.Server
		Expected string
	}{
		{
			Server: &highv3.Server{
				URL: "{GRAPHQL_SERVER_URL}",
			},
			Expected: serverURL,
		},
	}

	for _, tc := range testCases {
		parsedURL, err := parseServerURL(tc.Server, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Equal(t, parsedURL, tc.Expected)
	}
}
