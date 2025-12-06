package schema

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseSecuritySchemeType(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    SecuritySchemeType
		expectError bool
	}{
		{
			name:        "valid apiKey",
			input:       "apiKey",
			expected:    APIKeyScheme,
			expectError: false,
		},
		{
			name:        "valid http",
			input:       "http",
			expected:    HTTPAuthScheme,
			expectError: false,
		},
		{
			name:        "valid basic",
			input:       "basic",
			expected:    BasicAuthScheme,
			expectError: false,
		},
		{
			name:        "valid cookie",
			input:       "cookie",
			expected:    CookieAuthScheme,
			expectError: false,
		},
		{
			name:        "valid oauth2",
			input:       "oauth2",
			expected:    OAuth2Scheme,
			expectError: false,
		},
		{
			name:        "valid openIdConnect",
			input:       "openIdConnect",
			expected:    OpenIDConnectScheme,
			expectError: false,
		},
		{
			name:        "valid mutualTLS",
			input:       "mutualTLS",
			expected:    MutualTLSScheme,
			expectError: false,
		},
		{
			name:        "invalid scheme",
			input:       "invalid",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseSecuritySchemeType(tc.input)
			if tc.expectError {
				assert.ErrorIs(t, err, ErrInvalidSecuritySchemeType)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestSupportedSecuritySchemeTypes(t *testing.T) {
	schemes := SupportedSecuritySchemeTypes()
	assert.Assert(t, len(schemes) == 7)
	assert.Assert(t, contains(schemes, APIKeyScheme))
	assert.Assert(t, contains(schemes, HTTPAuthScheme))
	assert.Assert(t, contains(schemes, BasicAuthScheme))
	assert.Assert(t, contains(schemes, CookieAuthScheme))
	assert.Assert(t, contains(schemes, OAuth2Scheme))
	assert.Assert(t, contains(schemes, OpenIDConnectScheme))
	assert.Assert(t, contains(schemes, MutualTLSScheme))
}

func TestParseOAuthFlowType(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    OAuthFlowType
		expectError bool
	}{
		{
			name:        "valid authorizationCode",
			input:       "authorizationCode",
			expected:    AuthorizationCodeFlow,
			expectError: false,
		},
		{
			name:        "valid implicit",
			input:       "implicit",
			expected:    ImplicitFlow,
			expectError: false,
		},
		{
			name:        "valid password",
			input:       "password",
			expected:    PasswordFlow,
			expectError: false,
		},
		{
			name:        "valid clientCredentials",
			input:       "clientCredentials",
			expected:    ClientCredentialsFlow,
			expectError: false,
		},
		{
			name:        "invalid flow",
			input:       "invalid",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseOAuthFlowType(tc.input)
			if tc.expectError {
				assert.ErrorIs(t, err, ErrInvalidOAuthFlowType)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestSupportedParameterLocations(t *testing.T) {
	locations := SupportedParameterLocations()
	assert.Assert(t, len(locations) == 6)
	assert.Assert(t, containsParamLocation(locations, InQuery))
	assert.Assert(t, containsParamLocation(locations, InHeader))
	assert.Assert(t, containsParamLocation(locations, InPath))
	assert.Assert(t, containsParamLocation(locations, InCookie))
	assert.Assert(t, containsParamLocation(locations, InBody))
	assert.Assert(t, containsParamLocation(locations, InFormData))
}

// Helper function to check if a slice contains a SecuritySchemeType
func contains(slice []SecuritySchemeType, item SecuritySchemeType) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to check if a slice contains a ParameterLocation
func containsParamLocation(slice []ParameterLocation, item ParameterLocation) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
