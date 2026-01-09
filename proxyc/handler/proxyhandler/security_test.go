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

// TestNewOpenAPIv3Authenticator tests creating an authenticator
func TestNewOpenAPIv3Authenticator(t *testing.T) {
	t.Run("with_nil_components", func(t *testing.T) {
		doc := &highv3.Document{}
		auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, auth != nil)
		assert.Equal(t, 0, len(auth.securitySchemes))
	})

	t.Run("with_nil_security_schemes", func(t *testing.T) {
		doc := &highv3.Document{
			Components: &highv3.Components{},
		}
		auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, auth != nil)
		assert.Equal(t, 0, len(auth.securitySchemes))
	})

	t.Run("with_api_key_security_scheme", func(t *testing.T) {
		securitySchemes := orderedmap.New[string, *highv3.SecurityScheme]()
		securitySchemes.Set("api_key", &highv3.SecurityScheme{
			Type: "apiKey",
			Name: "X-API-Key",
			In:   "header",
		})

		doc := &highv3.Document{
			Components: &highv3.Components{
				SecuritySchemes: securitySchemes,
			},
		}

		auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, auth != nil)
		// API key without credentials won't create an authenticator
		assert.Equal(t, 0, len(auth.securitySchemes))
	})

	t.Run("with_http_basic_security_scheme", func(t *testing.T) {
		securitySchemes := orderedmap.New[string, *highv3.SecurityScheme]()
		securitySchemes.Set("basic_auth", &highv3.SecurityScheme{
			Type:   "http",
			Scheme: "basic",
		})

		doc := &highv3.Document{
			Components: &highv3.Components{
				SecuritySchemes: securitySchemes,
			},
		}

		auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, auth != nil)
		// Basic auth without credentials won't create an authenticator
		assert.Equal(t, 0, len(auth.securitySchemes))
	})

	t.Run("with_http_bearer_security_scheme", func(t *testing.T) {
		securitySchemes := orderedmap.New[string, *highv3.SecurityScheme]()
		securitySchemes.Set("bearer_auth", &highv3.SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		})

		doc := &highv3.Document{
			Components: &highv3.Components{
				SecuritySchemes: securitySchemes,
			},
		}

		auth, err := NewOpenAPIv3Authenticator(doc, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, auth != nil)
		// Bearer auth without credentials won't create an authenticator
		assert.Equal(t, 0, len(auth.securitySchemes))
	})
}

// TestGetAuthenticatorWithNoSchemes tests getting an authenticator with no schemes
func TestGetAuthenticatorWithNoSchemes(t *testing.T) {
	auth := &OpenAPIAuthenticator{
		securitySchemes: map[string]authscheme.HTTPClientAuthenticator{},
	}

	result := auth.GetAuthenticator(nil)
	assert.Assert(t, result == nil)
}
