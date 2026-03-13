package schema

import (
	"errors"
	"testing"

	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyResource_UnmarshalYAML_OpenAPI(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError error
		checkFunc   func(*testing.T, *RelixyResource)
	}{
		{
			name: "valid OpenAPI resource with inline spec",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api
  description: Test API description
definition:
  settings:
    headers: {}
    forwardHeaders:
      request: ["X-Hasura-Role", "Authorization", "X-Hasura-ddn-token"]
      response: []
  spec:
    openapi: "3.2.0"
    servers:
      - url: "{GRAPHQL_SERVER_URL}"
    paths:
      /v1/api/rest/artistbyname/{name}:
        get:
          operationId: artistByName
          x-rely-proxy-action:
            type: graphql
            request:
              query: |
                query artistByName($name: string!) { 
                  artist(where: {
                    name: { _eq: $name }
                  }) { 
                    name 
                  }
                }
              variables:
                name:
                  path: param.name
      /v1/api/rest/artists:
        get:
          operationId: artists
          x-rely-proxy-action:
            type: graphql
            request:
              query: |
                query artists($limit: Int = 10, $offset: Int = 0) { 
                  artist(limit: $limit, offset: $offset) { 
                    name 
                  }
                }
              variables:
                limit:
                  path: query.limit[0]
                  default:
                    value: 19
                offset:
                  path: query.offset[0]
                  default:
                    value: 0`,
			checkFunc: func(t *testing.T, r *RelixyResource) {
				base := r.GetBaseResource()
				assert.Equal(t, "v1", base.Version)
				assert.Equal(t, "test-api", base.Metadata.Name)
				assert.Equal(t, "Test API description", base.Metadata.Description)
			},
		},
		{
			name: "valid OpenAPI resource minimal metadata",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: minimal-api
definition:
  spec:
    openapi: "3.0.0"
    info:
      title: Minimal API
      version: "1.0.0"
    paths: {}`,
			checkFunc: func(t *testing.T, r *RelixyResource) {
				base := r.GetBaseResource()
				assert.Equal(t, "v1", base.Version)
				assert.Equal(t, "minimal-api", base.Metadata.Name)
				assert.Equal(t, "", base.Metadata.Description)
			},
		},
		{
			name: "valid OpenAPI resource with ref",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: ref-api
definition:
  ref: https://example.com/openapi.yaml`,
			checkFunc: func(t *testing.T, r *RelixyResource) {
				base := r.GetBaseResource()
				assert.Equal(t, "ref-api", base.Metadata.Name)
			},
		},
		{
			name: "missing definition",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api`,
			expectError: ErrRelixyResourceDefinitionRequired,
		},
		{
			name: "unsupported kind",
			yamlData: `version: v1
kind: Unknown
metadata:
  name: test-api
definition:
  spec:
    openapi: "3.0.0"
    info:
      title: Test
      version: "1.0.0"
    paths: {}`,
			expectError: ErrUnsupportedRelixyResourceKind,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var r RelixyResource
			err := yaml.Load([]byte(tc.yamlData), &r)

			if tc.expectError != nil {
				assert.Assert(t, errors.Is(err, tc.expectError), "expected error %v, got %v", tc.expectError, err)
				return
			}

			assert.NilError(t, err)
			assert.Assert(t, r.RelixyResource != nil)

			if tc.checkFunc != nil {
				tc.checkFunc(t, &r)
			}
		})
	}
}

func TestRelixyResource_UnmarshalYAML_RelyAuth(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError error
		checkFunc   func(*testing.T, *RelixyResource)
	}{
		{
			name: "valid RelyAuth resource with noAuth mode",
			yamlData: `version: v1
kind: RelyAuth
definition:
  modes:
    - mode: noAuth`,
			checkFunc: func(t *testing.T, r *RelixyResource) {
				base := r.GetBaseResource()
				assert.Equal(t, "v1", base.Version)
				assert.Equal(t, "RelyAuth", string(base.Kind))
			},
		},
		{
			name: "valid RelyAuth resource with apiKey mode",
			yamlData: `version: v1
kind: RelyAuth
definition:
  modes:
    - mode: apiKey
      tokenLocation:
        in: header
        name: Authorization
      value:
        value: secret-token
      sessionVariables: {}`,
			checkFunc: func(t *testing.T, r *RelixyResource) {
				base := r.GetBaseResource()
				assert.Equal(t, "RelyAuth", string(base.Kind))
			},
		},
		{
			name: "RelyAuth missing definition",
			yamlData: `version: v1
kind: RelyAuth
metadata:
  name: auth-resource`,
			expectError: ErrRelixyResourceDefinitionRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var r RelixyResource
			err := yaml.Load([]byte(tc.yamlData), &r)

			if tc.expectError != nil {
				assert.Assert(t, errors.Is(err, tc.expectError), "expected error %v, got %v", tc.expectError, err)
				return
			}

			assert.NilError(t, err)
			assert.Assert(t, r.RelixyResource != nil)

			if tc.checkFunc != nil {
				tc.checkFunc(t, &r)
			}
		})
	}
}
