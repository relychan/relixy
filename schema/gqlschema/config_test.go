package gqlschema

import (
	"encoding/json"
	"testing"

	"github.com/relychan/relixy/schema/base_schema"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyGraphQLConfig_JSONMarshal(t *testing.T) {
	config := RelixyGraphQLConfig{
		ScalarTypeMapping: map[string]base_schema.PrimitiveType{
			"DateTime": base_schema.String,
			"UUID":     base_schema.String,
			"Int64":    base_schema.Integer,
		},
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLConfig
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(result.ScalarTypeMapping))
	assert.Equal(t, base_schema.String, result.ScalarTypeMapping["DateTime"])
	assert.Equal(t, base_schema.String, result.ScalarTypeMapping["UUID"])
	assert.Equal(t, base_schema.Integer, result.ScalarTypeMapping["Int64"])
}

func TestRelixyGraphQLConfig_YAMLMarshal(t *testing.T) {
	config := RelixyGraphQLConfig{
		ScalarTypeMapping: map[string]base_schema.PrimitiveType{
			"DateTime": base_schema.String,
			"UUID":     base_schema.String,
		},
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLConfig
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(result.ScalarTypeMapping))
	assert.Equal(t, base_schema.String, result.ScalarTypeMapping["DateTime"])
	assert.Equal(t, base_schema.String, result.ScalarTypeMapping["UUID"])
}

func TestRelixyGraphQLConfig_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyGraphQLConfig)
	}{
		{
			name: "complete config with scalar type mapping",
			jsonData: `{
				"scalarTypeMapping": {
					"DateTime": "string",
					"UUID": "string",
					"Int64": "integer",
					"Float64": "number",
					"Boolean": "boolean"
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyGraphQLConfig) {
				assert.Equal(t, 5, len(config.ScalarTypeMapping))
				assert.Equal(t, base_schema.String, config.ScalarTypeMapping["DateTime"])
				assert.Equal(t, base_schema.String, config.ScalarTypeMapping["UUID"])
				assert.Equal(t, base_schema.Integer, config.ScalarTypeMapping["Int64"])
				assert.Equal(t, base_schema.Number, config.ScalarTypeMapping["Float64"])
				assert.Equal(t, base_schema.Boolean, config.ScalarTypeMapping["Boolean"])
			},
		},
		{
			name:        "empty config",
			jsonData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyGraphQLConfig) {
				assert.Assert(t, len(config.ScalarTypeMapping) == 0)
			},
		},
		{
			name: "config with null scalar type mapping",
			jsonData: `{
				"scalarTypeMapping": null
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyGraphQLConfig) {
				assert.Assert(t, config.ScalarTypeMapping == nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config RelixyGraphQLConfig
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

func TestRelixyGraphQLConfig_YAMLUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyGraphQLConfig)
	}{
		{
			name: "complete config with scalar type mapping",
			yamlData: `scalarTypeMapping:
  DateTime: string
  UUID: string
  Int64: integer
  Float64: number
  Boolean: boolean`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyGraphQLConfig) {
				assert.Equal(t, 5, len(config.ScalarTypeMapping))
				assert.Equal(t, base_schema.String, config.ScalarTypeMapping["DateTime"])
				assert.Equal(t, base_schema.String, config.ScalarTypeMapping["UUID"])
				assert.Equal(t, base_schema.Integer, config.ScalarTypeMapping["Int64"])
				assert.Equal(t, base_schema.Number, config.ScalarTypeMapping["Float64"])
				assert.Equal(t, base_schema.Boolean, config.ScalarTypeMapping["Boolean"])
			},
		},
		{
			name:        "empty config",
			yamlData:    `{}`,
			expectError: false,
			checkFunc: func(t *testing.T, config *RelixyGraphQLConfig) {
				assert.Assert(t, len(config.ScalarTypeMapping) == 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config RelixyGraphQLConfig
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
