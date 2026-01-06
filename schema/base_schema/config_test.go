package base_schema

import (
	"encoding/json"
	"testing"

	"github.com/relychan/gohttpc/loadbalancer"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyHealthCheckConfig_JSONMarshal(t *testing.T) {
	interval := 30
	timeout := 5
	config := RelixyHealthCheckConfig{
		HTTP: &loadbalancer.HTTPHealthCheckConfig{
			Path:     "/health",
			Interval: &interval,
			Timeout:  &timeout,
		},
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyHealthCheckConfig
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.HTTP != nil)
	assert.Equal(t, "/health", result.HTTP.Path)
}

func TestRelixyHealthCheckConfig_YAMLMarshal(t *testing.T) {
	interval := 30
	timeout := 5
	config := RelixyHealthCheckConfig{
		HTTP: &loadbalancer.HTTPHealthCheckConfig{
			Path:     "/health",
			Interval: &interval,
			Timeout:  &timeout,
		},
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyHealthCheckConfig
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.HTTP != nil)
	assert.Equal(t, "/health", result.HTTP.Path)
}

func TestRelixyHealthCheckConfig_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyHealthCheckConfig)
	}{
		{
			name: "complete config with HTTP health check",
			jsonData: `{
				"http": {
					"path": "/health",
					"interval": 30,
					"timeout": 5
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyHealthCheckConfig) {
				assert.Assert(t, config.HTTP != nil)
				assert.Equal(t, "/health", config.HTTP.Path)
				assert.Assert(t, config.HTTP.Interval != nil)
				assert.Assert(t, config.HTTP.Timeout != nil)
			},
		},
		{
			name:        "empty config",
			jsonData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyHealthCheckConfig) {
				assert.Assert(t, config.HTTP == nil)
			},
		},
		{
			name: "config with null HTTP",
			jsonData: `{
				"http": null
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyHealthCheckConfig) {
				assert.Assert(t, config.HTTP == nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config RelixyHealthCheckConfig
			err := json.Unmarshal([]byte(tc.jsonData), &config)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &config)
				}
			}
		})
	}
}

func TestRelixyHealthCheckConfig_YAMLUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyHealthCheckConfig)
	}{
		{
			name: "complete config with HTTP health check",
			yamlData: `http:
  path: /health
  interval: 30
  timeout: 5`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyHealthCheckConfig) {
				assert.Assert(t, config.HTTP != nil)
				assert.Equal(t, "/health", config.HTTP.Path)
				assert.Assert(t, config.HTTP.Interval != nil)
				assert.Assert(t, config.HTTP.Timeout != nil)
			},
		},
		{
			name:        "empty config",
			yamlData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyHealthCheckConfig) {
				assert.Assert(t, config.HTTP == nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config RelixyHealthCheckConfig
			err := yaml.Unmarshal([]byte(tc.yamlData), &config)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &config)
				}
			}
		})
	}
}

