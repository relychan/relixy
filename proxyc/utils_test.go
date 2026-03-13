package proxyc

import (
	"testing"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gotest.tools/v3/assert"
)

func TestParseServerURL(t *testing.T) {
	serverURL := "http://localhost:8080"
	t.Setenv("GRAPHQL_SERVER_URL", serverURL)
	t.Setenv("PORT", "8080")

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
		{
			Server: &highv3.Server{
				URL: "{GRAPHQL_SERVER_URL}/v1/graphql",
			},
			Expected: serverURL + "/v1/graphql",
		},
		{
			Server: &highv3.Server{
				URL: "http://{FOO}:{PORT}",
				Variables: func() *orderedmap.Map[string, *highv3.ServerVariable] {
					vars := orderedmap.New[string, *highv3.ServerVariable]()
					vars.Set("FOO", &highv3.ServerVariable{
						Default: "bar",
					})
					return vars
				}(),
			},
			Expected: "http://bar:8080",
		},
	}

	for _, tc := range testCases {
		parsedURL, err := parseServerURL(tc.Server, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Equal(t, parsedURL, tc.Expected)
	}
}
