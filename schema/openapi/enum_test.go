package openapi

import (
	"slices"
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
	assert.Assert(t, slices.Contains(schemes, APIKeyScheme))
	assert.Assert(t, slices.Contains(schemes, HTTPAuthScheme))
	assert.Assert(t, slices.Contains(schemes, BasicAuthScheme))
	assert.Assert(t, slices.Contains(schemes, CookieAuthScheme))
	assert.Assert(t, slices.Contains(schemes, OAuth2Scheme))
	assert.Assert(t, slices.Contains(schemes, OpenIDConnectScheme))
	assert.Assert(t, slices.Contains(schemes, MutualTLSScheme))
}
