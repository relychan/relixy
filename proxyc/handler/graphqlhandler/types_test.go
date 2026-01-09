package graphqlhandler

import (
	"testing"

	"gotest.tools/v3/assert"
)

// TestRelixyCustomGraphQLResponseConfigIsZero tests the IsZero method
func TestRelixyCustomGraphQLResponseConfigIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		config   RelixyCustomGraphQLResponseConfig
		expected bool
	}{
		{
			name:     "empty_config",
			config:   RelixyCustomGraphQLResponseConfig{},
			expected: true,
		},
		{
			name: "with_http_error_code",
			config: RelixyCustomGraphQLResponseConfig{
				HTTPErrorCode: func() *int { i := 400; return &i }(),
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestNewRelixyCustomGraphQLResponse tests creating custom GraphQL response
func TestNewRelixyCustomGraphQLResponse(t *testing.T) {
	t.Run("nil_config", func(t *testing.T) {
		result, err := NewRelixyCustomGraphQLResponse(nil, nil)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})

	t.Run("empty_config", func(t *testing.T) {
		config := &RelixyCustomGraphQLResponseConfig{}
		result, err := NewRelixyCustomGraphQLResponse(config, nil)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})

	t.Run("with_http_error_code", func(t *testing.T) {
		errorCode := 400
		config := &RelixyCustomGraphQLResponseConfig{
			HTTPErrorCode: &errorCode,
		}
		result, err := NewRelixyCustomGraphQLResponse(config, nil)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, 400, *result.HTTPErrorCode)
	})
}

// TestRelixyCustomGraphQLResponseIsZero tests the IsZero method
func TestRelixyCustomGraphQLResponseIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		response RelixyCustomGraphQLResponse
		expected bool
	}{
		{
			name:     "empty_response",
			response: RelixyCustomGraphQLResponse{},
			expected: true,
		},
		{
			name: "with_http_error_code",
			response: RelixyCustomGraphQLResponse{
				HTTPErrorCode: func() *int { i := 400; return &i }(),
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.response.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}
