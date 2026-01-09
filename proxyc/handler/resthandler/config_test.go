package resthandler

import (
	"testing"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/goutils"
	"gotest.tools/v3/assert"
)

// TestRelixyRESTRequestConfigIsZero tests the IsZero method
func TestRelixyRESTRequestConfigIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		config   RelixyRESTRequestConfig
		expected bool
	}{
		{
			name:     "empty_config",
			config:   RelixyRESTRequestConfig{},
			expected: true,
		},
		{
			name: "with_path",
			config: RelixyRESTRequestConfig{
				Path: "/test",
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

// TestRelixyCustomRESTResponseConfigIsZero tests the IsZero method
func TestRelixyCustomRESTResponseConfigIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		config   RelixyCustomRESTResponseConfig
		expected bool
	}{
		{
			name:     "empty_config",
			config:   RelixyCustomRESTResponseConfig{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestCustomRESTResponseIsZero tests the IsZero method for customRESTResponse
func TestCustomRESTResponseIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		response customRESTResponse
		expected bool
	}{
		{
			name:     "empty_response",
			response: customRESTResponse{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.response.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestCustomRESTRequestIsZero tests the IsZero method for customRESTRequest
func TestCustomRESTRequestIsZero(t *testing.T) {
	testCases := []struct {
		name     string
		request  customRESTRequest
		expected bool
	}{
		{
			name:     "empty_request",
			request:  customRESTRequest{},
			expected: true,
		},
		{
			name: "with_path",
			request: customRESTRequest{
				Path: "/test",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.request.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestNewCustomRESTResponse tests creating custom REST response
func TestNewCustomRESTResponse(t *testing.T) {
	t.Run("nil_config", func(t *testing.T) {
		result, err := newCustomRESTResponse(nil, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})

	t.Run("empty_config", func(t *testing.T) {
		config := &RelixyCustomRESTResponseConfig{}
		result, err := newCustomRESTResponse(config, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})
}

// TestNewCustomRESTRequestFromConfig tests creating custom REST request
func TestNewCustomRESTRequestFromConfig(t *testing.T) {
	t.Run("nil_config", func(t *testing.T) {
		result, err := newCustomRESTRequestFromConfig(nil, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})

	t.Run("empty_config", func(t *testing.T) {
		config := &RelixyRESTRequestConfig{}
		result, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result == nil)
	})

	t.Run("with_path", func(t *testing.T) {
		config := &RelixyRESTRequestConfig{
			Path: "/custom/path",
		}
		result, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, "/custom/path", result.Path)
	})

	t.Run("with_headers", func(t *testing.T) {
		config := &RelixyRESTRequestConfig{
			Headers: map[string]jmes.FieldMappingEntryStringConfig{
				"X-Custom-Header": {
					Path: goutils.ToPtr("headers.authorization"),
				},
			},
		}
		result, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Assert(t, len(result.Headers) > 0)
	})

	t.Run("with_path_and_headers", func(t *testing.T) {
		config := &RelixyRESTRequestConfig{
			Path: "/api/endpoint",
			Headers: map[string]jmes.FieldMappingEntryStringConfig{
				"Authorization": {
					Path: goutils.ToPtr("headers.auth"),
				},
			},
		}
		result, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, "/api/endpoint", result.Path)
		assert.Assert(t, len(result.Headers) > 0)
	})

	t.Run("invalid_headers_config", func(t *testing.T) {
		config := &RelixyRESTRequestConfig{
			Headers: map[string]jmes.FieldMappingEntryStringConfig{
				"X-Invalid": {},
			},
		}
		_, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		assert.Assert(t, err != nil)
		assert.ErrorContains(t, err, "failed to initialize custom request headers")
	})

	t.Run("with_body_only", func(t *testing.T) {
		bodyConfig := gotransform.TemplateTransformerConfig{}
		config := &RelixyRESTRequestConfig{
			Body: &bodyConfig,
		}
		_, err := newCustomRESTRequestFromConfig(config, goenvconf.GetOSEnv)
		// May fail due to invalid body config, which is expected
		if err != nil {
			assert.ErrorContains(t, err, "failed to initialize custom request body")
		}
	})
}
