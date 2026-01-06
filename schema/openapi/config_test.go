package openapi

import (
	"encoding/json"
	"testing"

	"github.com/hasura/goenvconf"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyOpenAPISettings_JSONMarshal(t *testing.T) {
	expose := true
	testKey := "test-key"
	config := RelixyOpenAPISettings{
		Expose:   &expose,
		BasePath: "/api/v1",
		Headers: map[string]goenvconf.EnvString{
			"X-API-Key": {Value: &testKey},
		},
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyOpenAPISettings
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.Expose != nil)
	assert.Equal(t, true, *result.Expose)
	assert.Equal(t, "/api/v1", result.BasePath)
	assert.Equal(t, 1, len(result.Headers))
}

func TestRelixyOpenAPISettings_YAMLMarshal(t *testing.T) {
	expose := false
	config := RelixyOpenAPISettings{
		Expose:   &expose,
		BasePath: "/api/v2",
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyOpenAPISettings
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.Expose != nil)
	assert.Equal(t, false, *result.Expose)
	assert.Equal(t, "/api/v2", result.BasePath)
}

func TestRelixyOpenAPISettings_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPISettings)
	}{
		{
			name: "complete settings",
			jsonData: `{
				"expose": true,
				"basePath": "/api/v1",
				"forwardHeaders": {
					"request": ["Authorization", "X-Request-ID"],
					"response": ["X-Response-ID"]
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, settings *RelixyOpenAPISettings) {
				assert.Assert(t, settings.Expose != nil)
				assert.Equal(t, true, *settings.Expose)
				assert.Equal(t, "/api/v1", settings.BasePath)
				assert.Assert(t, settings.ForwardHeaders != nil)
				assert.Equal(t, 2, len(settings.ForwardHeaders.Request))
				assert.Equal(t, 1, len(settings.ForwardHeaders.Response))
			},
		},
		{
			name: "settings with health check",
			jsonData: `{
				"basePath": "/api",
				"healthCheck": {
					"http": {
						"path": "/health",
						"interval": 30,
						"timeout": 5
					}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, settings *RelixyOpenAPISettings) {
				assert.Equal(t, "/api", settings.BasePath)
				assert.Assert(t, settings.HealthCheck != nil)
				assert.Assert(t, settings.HealthCheck.HTTP != nil)
				assert.Equal(t, "/health", settings.HealthCheck.HTTP.Path)
			},
		},
		{
			name:        "empty settings",
			jsonData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, settings *RelixyOpenAPISettings) {
				assert.Assert(t, settings.Expose == nil)
				assert.Equal(t, "", settings.BasePath)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var settings RelixyOpenAPISettings
			err := json.Unmarshal([]byte(tc.jsonData), &settings)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &settings)
				}
			}
		})
	}
}

func TestRelixyOpenAPISettings_YAMLUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPISettings)
	}{
		{
			name: "complete settings",
			yamlData: `expose: true
basePath: /api/v1
forwardHeaders:
  request:
    - Authorization
    - X-Request-ID
  response:
    - X-Response-ID`,
			expectError: false,
			checkFunc: func(t *testing.T, settings *RelixyOpenAPISettings) {
				assert.Assert(t, settings.Expose != nil)
				assert.Equal(t, true, *settings.Expose)
				assert.Equal(t, "/api/v1", settings.BasePath)
				assert.Assert(t, settings.ForwardHeaders != nil)
				assert.Equal(t, 2, len(settings.ForwardHeaders.Request))
				assert.Equal(t, 1, len(settings.ForwardHeaders.Response))
			},
		},
		{
			name:        "empty settings",
			yamlData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, settings *RelixyOpenAPISettings) {
				assert.Assert(t, settings.Expose == nil)
				assert.Equal(t, "", settings.BasePath)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var settings RelixyOpenAPISettings
			err := yaml.Unmarshal([]byte(tc.yamlData), &settings)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &settings)
				}
			}
		})
	}
}

func TestRelixyOpenAPIForwardHeadersConfig_JSONMarshal(t *testing.T) {
	config := RelixyOpenAPIForwardHeadersConfig{
		Request:  []string{"Authorization", "X-Request-ID"},
		Response: []string{"X-Response-ID", "X-Trace-ID"},
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyOpenAPIForwardHeadersConfig
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(result.Request))
	assert.Equal(t, 2, len(result.Response))
	assert.Equal(t, "Authorization", result.Request[0])
	assert.Equal(t, "X-Response-ID", result.Response[0])
}

func TestRelixyOpenAPIForwardHeadersConfig_YAMLMarshal(t *testing.T) {
	config := RelixyOpenAPIForwardHeadersConfig{
		Request:  []string{"Authorization"},
		Response: []string{"X-Response-ID"},
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyOpenAPIForwardHeadersConfig
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(result.Request))
	assert.Equal(t, 1, len(result.Response))
}


