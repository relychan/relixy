package proxyhandler

import (
	"testing"

	"github.com/hasura/goenvconf"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/relychan/gohttpc/authc/authscheme"
	"github.com/relychan/relixy/schema/openapi"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestNewOpenAPIv3Authenticator_EmptyDocument(t *testing.T) {
	doc := &highv3.Document{}

	auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Assert(t, auth != nil)
	assert.Equal(t, 0, len(auth.securitySchemes))
}

func TestNewOpenAPIv3Authenticator_NoSecuritySchemes(t *testing.T) {
	doc := &highv3.Document{
		Components: &highv3.Components{},
	}

	auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Assert(t, auth != nil)
	assert.Equal(t, 0, len(auth.securitySchemes))
}

func TestNewOpenAPIv3Authenticator_WithAPIKey(t *testing.T) {
	apiKeyValue := goenvconf.NewEnvStringValue("test-api-key")
	creds := APIKeyCredentials{
		APIKey: &apiKeyValue,
	}

	credsData, err := yaml.Marshal(creds)
	assert.NilError(t, err)

	var credsNode yaml.Node
	err = yaml.Unmarshal(credsData, &credsNode)
	assert.NilError(t, err)

	extensions := orderedmap.New[string, *yaml.Node]()
	extensions.Set(openapi.XRelySecurityCredentials, &credsNode)

	securitySchemes := orderedmap.New[string, *highv3.SecurityScheme]()
	securitySchemes.Set("apiKey", &highv3.SecurityScheme{
		Type:       string(openapi.APIKeyScheme),
		Name:       "X-API-Key",
		In:         "header",
		Extensions: extensions,
	})

	doc := &highv3.Document{
		Components: &highv3.Components{
			SecuritySchemes: securitySchemes,
		},
	}

	auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Assert(t, auth != nil)
	assert.Equal(t, 1, len(auth.securitySchemes))
}

func TestNewOpenAPIv3Authenticator_WithBasicAuth(t *testing.T) {
	username := goenvconf.NewEnvStringValue("test-user")
	password := goenvconf.NewEnvStringValue("test-pass")
	creds := BasicCredentials{
		Username: &username,
		Password: &password,
	}

	credsData, err := yaml.Marshal(creds)
	assert.NilError(t, err)

	var credsNode yaml.Node
	err = yaml.Unmarshal(credsData, &credsNode)
	assert.NilError(t, err)

	extensions := orderedmap.New[string, *yaml.Node]()
	extensions.Set(openapi.XRelySecurityCredentials, &credsNode)

	securitySchemes := orderedmap.New[string, *highv3.SecurityScheme]()
	securitySchemes.Set("basicAuth", &highv3.SecurityScheme{
		Type:       string(openapi.HTTPAuthScheme),
		Scheme:     string(openapi.BasicAuthScheme),
		Name:       "Authorization",
		Extensions: extensions,
	})

	doc := &highv3.Document{
		Components: &highv3.Components{
			SecuritySchemes: securitySchemes,
		},
	}

	auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Assert(t, auth != nil)
	assert.Equal(t, 1, len(auth.securitySchemes))
}

func TestOpenAPIAuthenticator_GetAuthenticator_Empty(t *testing.T) {
	auth := &OpenAPIAuthenticator{
		securitySchemes: make(map[string]authscheme.HTTPClientAuthenticator),
	}

	result := auth.GetAuthenticator(nil)
	assert.Assert(t, result == nil)
}

func TestOpenAPIAuthenticator_GetAuthenticator_WithOptionalSecurity(t *testing.T) {
	auth := &OpenAPIAuthenticator{
		securitySchemes: make(map[string]authscheme.HTTPClientAuthenticator),
		security: []*base.SecurityRequirement{
			{
				ContainsEmptyRequirement: true,
			},
		},
	}

	result := auth.GetAuthenticator(nil)
	assert.Assert(t, result == nil)
}

func TestParseStringFromExtensions(t *testing.T) {
	testCases := []struct {
		name          string
		key           string
		value         string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "valid string",
			key:           "test-key",
			value:         "test-value",
			expectedValue: "test-value",
			expectError:   false,
		},
		{
			name:          "empty string",
			key:           "test-key",
			value:         "",
			expectedValue: "",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			extensions := orderedmap.New[string, *yaml.Node]()

			var valueNode yaml.Node
			err := yaml.Unmarshal([]byte(tc.value), &valueNode)
			assert.NilError(t, err)

			extensions.Set(tc.key, &valueNode)

			result, err := parseStringFromExtensions(extensions, tc.key)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, tc.expectedValue, result)
			}
		})
	}
}
