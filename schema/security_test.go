package schema

import (
	"encoding/json"
	"testing"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gohttpc/authc/authscheme"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestAPIKeyAuthConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		config      *APIKeyAuthConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &APIKeyAuthConfig{
				Type:  APIKeyScheme,
				Name:  "X-API-Key",
				In:    authscheme.InHeader,
				Value: goenvconf.NewEnvStringValue("test-key"),
			},
			expectError: false,
		},
		{
			name: "missing name",
			config: &APIKeyAuthConfig{
				Type:  APIKeyScheme,
				In:    authscheme.InHeader,
				Value: goenvconf.NewEnvStringValue("test-key"),
			},
			expectError: true,
		},
		{
			name: "invalid location",
			config: &APIKeyAuthConfig{
				Type:  APIKeyScheme,
				Name:  "X-API-Key",
				In:    authscheme.AuthLocation("invalid"),
				Value: goenvconf.NewEnvStringValue("test-key"),
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestAPIKeyAuthConfig_GetType(t *testing.T) {
	config := &APIKeyAuthConfig{}
	assert.Equal(t, APIKeyScheme, config.GetType())
}

func TestHTTPAuthConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		config      *HTTPAuthConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &HTTPAuthConfig{
				Type:   HTTPAuthScheme,
				Scheme: "bearer",
				Header: "Authorization",
				Value:  goenvconf.NewEnvStringValue("token"),
			},
			expectError: false,
		},
		{
			name: "missing scheme",
			config: &HTTPAuthConfig{
				Type:   HTTPAuthScheme,
				Header: "Authorization",
				Value:  goenvconf.NewEnvStringValue("token"),
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestHTTPAuthConfig_GetType(t *testing.T) {
	config := &HTTPAuthConfig{}
	assert.Equal(t, HTTPAuthScheme, config.GetType())
}

func TestBasicAuthConfig_Validate(t *testing.T) {
	config := &BasicAuthConfig{
		Type:     BasicAuthScheme,
		Username: goenvconf.NewEnvStringValue("user"),
		Password: goenvconf.NewEnvStringValue("pass"),
	}
	err := config.Validate()
	assert.NilError(t, err)
}

func TestBasicAuthConfig_GetType(t *testing.T) {
	config := &BasicAuthConfig{}
	assert.Equal(t, BasicAuthScheme, config.GetType())
}

func TestOAuth2Config_Validate(t *testing.T) {
	tokenURL := goenvconf.NewEnvStringValue("https://example.com/token")
	clientID := goenvconf.NewEnvStringValue("client-id")
	clientSecret := goenvconf.NewEnvStringValue("client-secret")

	testCases := []struct {
		name        string
		config      *OAuth2Config
		expectError bool
	}{
		{
			name: "valid client credentials flow",
			config: &OAuth2Config{
				Type: OAuth2Scheme,
				Flows: map[OAuthFlowType]OAuthFlow{
					ClientCredentialsFlow: {
						TokenURL:     &tokenURL,
						ClientID:     &clientID,
						ClientSecret: &clientSecret,
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing client ID for client credentials",
			config: &OAuth2Config{
				Type: OAuth2Scheme,
				Flows: map[OAuthFlowType]OAuthFlow{
					ClientCredentialsFlow: {
						TokenURL:     &tokenURL,
						ClientSecret: &clientSecret,
					},
				},
			},
			expectError: true,
		},
		{
			name: "missing flows",
			config: &OAuth2Config{
				Type:  OAuth2Scheme,
				Flows: map[OAuthFlowType]OAuthFlow{},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestOAuth2Config_GetType(t *testing.T) {
	config := &OAuth2Config{Type: OAuth2Scheme}
	assert.Equal(t, OAuth2Scheme, config.GetType())
}

func TestOpenIDConnectConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		config      *OpenIDConnectConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &OpenIDConnectConfig{
				Type:             OpenIDConnectScheme,
				OpenIDConnectURL: "https://example.com/.well-known/openid-configuration",
			},
			expectError: false,
		},
		{
			name: "missing URL",
			config: &OpenIDConnectConfig{
				Type: OpenIDConnectScheme,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestOpenIDConnectConfig_GetType(t *testing.T) {
	config := &OpenIDConnectConfig{Type: OpenIDConnectScheme}
	assert.Equal(t, OpenIDConnectScheme, config.GetType())
}

func TestCookieAuthConfig_Validate(t *testing.T) {
	config := &CookieAuthConfig{Type: CookieAuthScheme}
	err := config.Validate()
	assert.NilError(t, err)
}

func TestCookieAuthConfig_GetType(t *testing.T) {
	config := &CookieAuthConfig{}
	assert.Equal(t, CookieAuthScheme, config.GetType())
}

func TestMutualTLSAuthConfig_Validate(t *testing.T) {
	config := &MutualTLSAuthConfig{Type: MutualTLSScheme}
	err := config.Validate()
	assert.NilError(t, err)
}

func TestMutualTLSAuthConfig_GetType(t *testing.T) {
	config := &MutualTLSAuthConfig{}
	assert.Equal(t, MutualTLSScheme, config.GetType())
}

func TestRelyProxySecurityScheme_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkType   SecuritySchemeType
	}{
		{
			name:        "valid apiKey",
			jsonData:    `{"type":"apiKey","name":"X-API-Key","in":"header","value":{"value":"test"}}`,
			expectError: false,
			checkType:   APIKeyScheme,
		},
		{
			name:        "valid http",
			jsonData:    `{"type":"http","scheme":"bearer","header":"Authorization","value":{"value":"token"}}`,
			expectError: false,
			checkType:   HTTPAuthScheme,
		},
		{
			name:        "valid basic",
			jsonData:    `{"type":"basic","username":{"value":"user"},"password":{"value":"pass"}}`,
			expectError: false,
			checkType:   BasicAuthScheme,
		},
		{
			name:        "valid cookie",
			jsonData:    `{"type":"cookie"}`,
			expectError: false,
			checkType:   CookieAuthScheme,
		},
		{
			name:        "valid mutualTLS",
			jsonData:    `{"type":"mutualTLS"}`,
			expectError: false,
			checkType:   MutualTLSScheme,
		},
		{
			name:        "invalid type",
			jsonData:    `{"type":"invalid"}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var scheme RelyProxySecurityScheme
			err := json.Unmarshal([]byte(tc.jsonData), &scheme)
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, tc.checkType, scheme.GetType())
			}
		})
	}
}

func TestRelyProxySecurityScheme_MarshalJSON(t *testing.T) {
	scheme := RelyProxySecurityScheme{
		SecuritySchemer: &CookieAuthConfig{
			Type: CookieAuthScheme,
		},
	}

	data, err := json.Marshal(scheme)
	assert.NilError(t, err)
	assert.Assert(t, len(data) > 0)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, "cookie", result["type"])
}

func TestRelyProxySecurityScheme_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkType   SecuritySchemeType
	}{
		{
			name: "valid apiKey",
			yamlData: `type: apiKey
name: X-API-Key
in: header
value:
  value: test`,
			expectError: false,
			checkType:   APIKeyScheme,
		},
		{
			name:        "valid cookie",
			yamlData:    `type: cookie`,
			expectError: false,
			checkType:   CookieAuthScheme,
		},
		{
			name:        "invalid type",
			yamlData:    `type: invalid`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var scheme RelyProxySecurityScheme
			err := yaml.Unmarshal([]byte(tc.yamlData), &scheme)
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, tc.checkType, scheme.GetType())
			}
		})
	}
}

func TestRelyProxySecurityScheme_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		scheme      RelyProxySecurityScheme
		expectError bool
	}{
		{
			name: "valid scheme",
			scheme: RelyProxySecurityScheme{
				SecuritySchemer: &CookieAuthConfig{
					Type: CookieAuthScheme,
				},
			},
			expectError: false,
		},
		{
			name: "nil schemer",
			scheme: RelyProxySecurityScheme{
				SecuritySchemer: nil,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.scheme.Validate()
			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
